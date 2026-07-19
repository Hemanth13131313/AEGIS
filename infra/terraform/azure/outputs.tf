output "aks_cluster_name" {
  description = "The name of the AKS cluster"
  value       = azurerm_kubernetes_cluster.aegis_aks.name
}

output "aks_kube_config" {
  description = "Kubeconfig for the AKS cluster"
  value       = azurerm_kubernetes_cluster.aegis_aks.kube_config_raw
  sensitive   = true
}

output "db_fqdn" {
  description = "The FQDN of the PostgreSQL database"
  value       = azurerm_postgresql_flexible_server.aegis_db_server.fqdn
}

output "redis_hostname" {
  description = "The hostname of the Redis cache"
  value       = azurerm_redis_cache.aegis_redis.hostname
}

output "eventhub_namespace_name" {
  description = "The name of the EventHub namespace"
  value       = azurerm_eventhub_namespace.aegis_eh_ns.name
}

output "acr_login_server" {
  description = "The login server for the Azure Container Registry"
  value       = azurerm_container_registry.aegis_acr.login_server
}
