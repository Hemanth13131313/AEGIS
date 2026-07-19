variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The region to deploy resources to"
  type        = string
  default     = "us-central1"
}

variable "environment" {
  description = "The environment name (e.g., production, staging)"
  type        = string
  default     = "production"
}

variable "gke_cluster_name" {
  description = "The name of the GKE cluster"
  type        = string
  default     = "aegis-cluster"
}

variable "db_tier" {
  description = "The machine tier for Cloud SQL PostgreSQL"
  type        = string
  default     = "db-custom-2-7680"
}

variable "redis_memory_size_gb" {
  description = "The memory size in GB for Redis"
  type        = number
  default     = 1
}
