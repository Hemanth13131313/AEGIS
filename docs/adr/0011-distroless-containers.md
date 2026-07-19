# ADR 0011: Use distroless container images for all production services

## Status
Accepted

## Context
Standard base images include shells, package managers, and various OS utilities that expand the attack surface and can be leveraged by attackers if they achieve remote code execution.

## Decision
All Go microservices will use `gcr.io/distroless/static-debian12`.
All Python microservices will use `gcr.io/distroless/python3-debian12`.

## Consequences
- No shell for debugging inside the container (use ephemeral debug containers instead).
- Smaller image sizes leading to faster pull times.
- Reduced risk of OS-level CVEs from shells and utilities.
