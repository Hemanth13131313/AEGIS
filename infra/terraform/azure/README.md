# AEGIS Azure Infrastructure

This module provisions the Azure infrastructure for the AEGIS platform.

## Prerequisites
- Azure Subscription
- `az` CLI installed and authenticated (`az login`)
- Terraform >= 1.8

## Usage

1. Initialize Terraform:
   ```bash
   terraform init -backend-config="resource_group_name=YOUR_RG" -backend-config="storage_account_name=YOUR_SA" -backend-config="container_name=tfstate" -backend-config="key=terraform.tfstate"
   ```

2. Plan the deployment:
   ```bash
   terraform plan -var="subscription_id=YOUR_SUBSCRIPTION_ID" -var="environment=production"
   ```

3. Apply the deployment:
   ```bash
   terraform apply -var="subscription_id=YOUR_SUBSCRIPTION_ID" -var="environment=production"
   ```

## Resources Created
| Resource | Description |
|----------|-------------|
| AKS Cluster | Managed Kubernetes cluster with Workload Identity and RBAC |
| PostgreSQL Flexible Server | Managed database instance (Zone-redundant) |
| Redis Cache | Standard tier Redis cache with TLS only |
| Event Hubs | Kafka-compatible event streaming endpoints |
| Container Registry | Azure Container Registry |
| Key Vault | Secret placeholders |
| Virtual Network & Subnets | Private network for infrastructure |
