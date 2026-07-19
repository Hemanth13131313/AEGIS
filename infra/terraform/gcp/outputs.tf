output "gke_cluster_endpoint" {
  description = "The IP address of the GKE cluster control plane"
  value       = google_container_cluster.aegis_cluster.endpoint
}

output "gke_cluster_name" {
  description = "The name of the GKE cluster"
  value       = google_container_cluster.aegis_cluster.name
}

output "db_connection_name" {
  description = "The connection name of the Cloud SQL instance"
  value       = google_sql_database_instance.aegis_db_instance.connection_name
}

output "redis_host" {
  description = "The IP address of the Redis instance"
  value       = google_redis_instance.aegis_redis.host
}

output "artifact_registry_url" {
  description = "The URL of the Artifact Registry repository"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.aegis_registry.repository_id}"
}
