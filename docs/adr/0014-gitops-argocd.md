# 0014 - GitOps with ArgoCD

## Status
Accepted

## Context
Deploying AEGIS via standard CI/CD pipelines can lead to configuration drift and requires giving the CI server broad access to the Kubernetes cluster.

## Decision
- We adopt ArgoCD as the primary GitOps controller for deploying AEGIS workloads.
- We configure ArgoCD with automated sync and self-healing.
- Flux kustomization is provided as a supported alternative.

## Consequences
- Requires installing the ArgoCD cluster addon.
- Enables continuous drift detection and automated rollback.
- Simplifies multi-environment promotion patterns.
