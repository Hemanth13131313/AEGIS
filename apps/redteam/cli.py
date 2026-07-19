#!/usr/bin/env python3
"""
AEGIS Red Team CLI.
"""
import asyncio
import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path

import click
from rich.console import Console
from rich.table import Table
from rich.progress import Progress, SpinnerColumn, TimeElapsedColumn
from rich import print as rprint
import jsonschema

from generators.base import TestCase
from generators.prompt_injection import PromptInjectionGenerator
from generators.jailbreak import JailbreakGenerator
from generators.data_exfiltration import DataExfiltrationGenerator
from generators.excessive_agency import ExcessiveAgencyGenerator
from generators.benign import BenignGenerator
from runner import RedTeamRunner, RunConfig, RunResult
from runner.result_store import ResultStore

console = Console()

@click.group()
@click.version_option()
def cli():
    pass

@cli.command('run')
@click.option('--target', required=True, help='AEGIS /scan URL')
@click.option('--api-key', envvar='AEGIS_API_KEY', default='')
@click.option('--cases', type=click.Path(exists=True), default='testcases/registry.json')
@click.option('--owasp', help='Filter by OWASP category')
@click.option('--concurrency', type=int, default=5)
@click.option('--timeout', type=float, default=10.0)
@click.option('--dry-run', is_flag=True)
@click.option('--fail-on-disagreement', is_flag=True)
@click.option('--output', type=click.Choice(['table', 'json', 'ci']), default='table')
def run(target, api_key, cases, owasp, concurrency, timeout, dry_run, fail_on_disagreement, output):
    """Run red team tests."""
    with open(cases, 'r') as f:
        cases_data = json.load(f)
        
    test_cases = [TestCase(**c) for c in cases_data]
    if owasp:
        test_cases = [c for c in test_cases if c.owasp_category == owasp]
        
    config = RunConfig(
        target_url=target,
        api_key=api_key,
        concurrency=concurrency,
        timeout_seconds=timeout,
        fail_on_disagreement=fail_on_disagreement,
        dry_run=dry_run
    )
    
    runner = RedTeamRunner(config)
    
    start_time = datetime.now(timezone.utc)
    results = []
    
    if output == 'table':
        with Progress(SpinnerColumn(), *Progress.get_default_columns(), TimeElapsedColumn(), console=console) as progress:
            task = progress.add_task("[cyan]Running red team tests...", total=len(test_cases))
            async def run_and_update():
                semaphore = asyncio.Semaphore(config.concurrency)
                import httpx
                async with httpx.AsyncClient(timeout=config.timeout_seconds) as client:
                    tasks = []
                    for c in test_cases:
                        async def wrap(c=c):
                            res = await runner.run_one(c, semaphore, client)
                            progress.advance(task)
                            return res
                        tasks.append(wrap())
                    return await asyncio.gather(*tasks)
            results = asyncio.run(run_and_update())
    else:
        results = asyncio.run(runner.run_all(test_cases))
        
    end_time = datetime.now(timezone.utc)
    
    # Sort for output
    results = sorted(results, key=lambda r: (r.passed, r.test_case_id))
    summary = runner.summary(results, start_time, end_time)
    
    # Save results
    results_dir = Path("results")
    store = ResultStore(results_dir, os.getenv("AEGIS_CLICKHOUSE_ADDR"))
    asyncio.run(store.save(results, summary))
    
    if output == 'table':
        table = Table(title="Red Team Results")
        table.add_column("Result")
        table.add_column("ID")
        table.add_column("Name")
        table.add_column("OWASP")
        table.add_column("ATLAS")
        table.add_column("Action Received")
        table.add_column("Expected")
        table.add_column("Latency")
        
        for r in results:
            if r.error:
                icon = "💥 Error"
                color = "red"
            elif r.passed:
                if r.disagreement:
                    icon = "⚠️ Pass (Diag)"
                    color = "yellow"
                else:
                    icon = "✅ Pass"
                    color = "green"
            else:
                icon = "❌ Fail"
                color = "red"
                
            table.add_row(
                f"[{color}]{icon}[/{color}]",
                r.test_case_id,
                r.test_case_name,
                r.owasp_category,
                r.atlas_technique,
                r.action_received,
                r.expected_action,
                f"{r.latency_ms:.0f}ms"
            )
            
        console.print(table)
        
        console.print(f"\n[bold]Summary[/bold]")
        console.print(f"Total: {summary.total} | Passed: {summary.passed} | Failed: {summary.failed} | Disagreements: {summary.disagreements}")
        color = "green" if summary.pass_rate >= 0.9 else "yellow" if summary.pass_rate >= 0.7 else "red"
        console.print(f"Pass Rate: [{color}]{summary.pass_rate*100:.1f}%[/{color}]")
        
    elif output == 'json' or output == 'ci':
        # JSON output
        pass # The json file is saved anyway
        
    if summary.failed > 0 or summary.errors > 0:
        sys.exit(1)

@cli.command('list-cases')
@click.option('--owasp', help='Filter by OWASP category')
@click.option('--atlas', help='Filter by ATLAS technique')
@click.option('--cases', type=click.Path(exists=True), default='testcases/registry.json')
def list_cases(owasp, atlas, cases):
    """List test cases from registry."""
    with open(cases, 'r') as f:
        cases_data = json.load(f)
        
    table = Table(title="Test Cases")
    table.add_column("ID")
    table.add_column("Name")
    table.add_column("OWASP")
    table.add_column("ATLAS")
    table.add_column("Expected Action")
    
    for c in cases_data:
        if owasp and c.get('owasp_category') != owasp:
            continue
        if atlas and c.get('atlas_technique') != atlas:
            continue
            
        table.add_row(
            c['id'],
            c['name'],
            c['owasp_category'],
            c['atlas_technique'],
            c['expected_action']
        )
        
    console.print(table)

@cli.command('validate')
@click.option('--cases', type=click.Path(exists=True), default='testcases/registry.json')
@click.option('--schema', type=click.Path(exists=True), default='testcases/schema.json')
def validate(cases, schema):
    """Validate test case registry against schema."""
    with open(cases, 'r') as f:
        cases_data = json.load(f)
    with open(schema, 'r') as f:
        schema_data = json.load(f)
        
    try:
        jsonschema.validate(instance=cases_data, schema=schema_data)
        console.print("[green]Registry is valid.[/green]")
        sys.exit(0)
    except jsonschema.exceptions.ValidationError as e:
        console.print(f"[red]Validation error: {e.message}[/red]")
        sys.exit(1)

@cli.command('generate')
@click.option('--type', 'gen_type', type=click.Choice(['prompt_injection', 'jailbreak', 'data_exfiltration', 'excessive_agency', 'benign']), required=True)
@click.option('--count', type=int, default=5)
def generate(gen_type, count):
    """Generate mock test cases."""
    if gen_type == 'prompt_injection':
        gen = PromptInjectionGenerator()
    elif gen_type == 'jailbreak':
        gen = JailbreakGenerator()
    elif gen_type == 'data_exfiltration':
        gen = DataExfiltrationGenerator()
    elif gen_type == 'excessive_agency':
        gen = ExcessiveAgencyGenerator()
    elif gen_type == 'benign':
        gen = BenignGenerator()
        
    cases = gen.generate(count)
    from dataclasses import asdict
    print(json.dumps([asdict(c) for c in cases], indent=2))

if __name__ == '__main__':
    cli()
