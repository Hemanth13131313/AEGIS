# Contributing to AEGIS

Thank you for your interest in contributing!

## Code of Conduct
By participating in this project, you agree to abide by our Code of Conduct. Please treat all contributors with respect.

## Development Setup
1. Clone the repository.
2. Copy `.env.dev.example` to `.env`.
3. Run `docker compose -f docker-compose.dev.yml up -d` to start the local environment.

## Branch Naming
Please use the following conventions:
- `feature/phase-N-short-desc`
- `fix/short-desc`

## PR Checklist
- [ ] Code is properly formatted (`gofmt`, `black`, `prettier`).
- [ ] No hardcoded secrets or credentials.
- [ ] Unit tests added/updated and passing.
- [ ] Documentation updated if necessary.
- [ ] Commit messages follow the conventional format.

## Commit Message Format
We use [Conventional Commits](https://www.conventionalcommits.org/):
```
type(scope): summary
```
Examples: `feat(gateway): add retry logic`, `fix(scanner): resolve timeout issue`.

## Testing Requirements
- Unit tests must be written for all new functionality.
- Integration tests should cover cross-service boundaries.
- Run `make test` locally before pushing.
