variable "subscription_id" {
  description = "The Azure subscription ID"
  type        = string
}

variable "resource_group_location" {
  description = "The location for the resource group"
  type        = string
  default     = "eastus"
}

variable "environment" {
  description = "The environment name"
  type        = string
  default     = "production"
}

variable "aks_cluster_name" {
  description = "The name of the AKS cluster"
  type        = string
  default     = "aegis-cluster"
}

variable "db_sku_name" {
  description = "The SKU Name for the PostgreSQL Flexible Server"
  type        = string
  default     = "Standard_D2ds_v4"
}
