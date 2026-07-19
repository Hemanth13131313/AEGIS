#!/usr/bin/env python3
"""Validates the ATLAS technique registry JSON against its schema."""
import json, sys
from pathlib import Path

def validate():
    data = json.loads(Path("atlas-techniques.json").read_text())
    techniques = data.get("techniques", [])
    ids = set()
    errors = []
    for t in techniques:
        if t["id"] in ids:
            errors.append(f"Duplicate ID: {t['id']}")
        ids.add(t["id"])
        if not t["id"].startswith("AML.T"):
            errors.append(f"Invalid ID format: {t['id']}")
    print(f"Validated {len(techniques)} techniques. Errors: {len(errors)}")
    for e in errors:
        print(f"  ERROR: {e}")
    if errors:
        sys.exit(1)

if __name__ == "__main__":
    validate()
