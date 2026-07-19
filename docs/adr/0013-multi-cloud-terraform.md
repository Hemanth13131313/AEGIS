# 0013 - Multi-Cloud Terraform

## Status
Accepted

## Context
AEGIS needs to run across different public clouds (AWS, GCP, Azure) without locking in to a single vendor. However, we want to share as much Kubernetes deployment logic as possible via Helm.

## Decision
- We will maintain one complete Terraform module per cloud (AWS, GCP, Azure) to handle managed services like Kubernetes (EKS/GKE/AKS), Databases, Caching, and Streaming natively.
- We will use shared Helm charts for the Kubernetes workloads, supplemented by `values-{cloud}.yaml` overrides for cloud-specific bindings like workload identities.

## Consequences
- Requires maintaining 3x infrastructure code surface.
- Justified by the strict multi-cloud requirement in the product spec.
- Makes the deployment highly adaptable to client environments.
