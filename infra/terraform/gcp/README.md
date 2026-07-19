# AEGIS GCP Infrastructure

This module provisions the GCP infrastructure for the AEGIS platform.

## Prerequisites
- GCP Project ID
- `gcloud` CLI installed and authenticated
- Terraform >= 1.8

## Usage

1. Initialize Terraform:
   ```bash
   terraform init -backend-config="bucket=YOUR_BUCKET" -backend-config="prefix=terraform/state"
   ```

2. Plan the deployment:
   ```bash
   terraform plan -var="project_id=YOUR_PROJECT_ID" -var="environment=production"
   ```

3. Apply the deployment:
   ```bash
   terraform apply -var="project_id=YOUR_PROJECT_ID" -var="environment=production"
   ```

## Resources Created
| Resource | Description |
|----------|-------------|
| GKE Autopilot Cluster | Managed Kubernetes cluster with Workload Identity enabled |
| Cloud SQL (PostgreSQL 16) | Managed database instance |
| Memorystore for Redis | High availability Redis cache with TLS |
| Cloud Pub/Sub | Kafka-compatible event streaming topics and subscriptions |
| Artifact Registry | Docker container repository |
| Secret Manager | Secret placeholders |
| VPC & Subnets | Private network for infrastructure |
