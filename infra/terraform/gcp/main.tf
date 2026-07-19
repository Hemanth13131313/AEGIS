provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_compute_network" "aegis_vpc" {
  name                    = "aegis-vpc-${var.environment}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "aegis_subnet" {
  name          = "aegis-subnet-${var.environment}"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.region
  network       = google_compute_network.aegis_vpc.id
  private_ip_google_access = true
}

resource "google_container_cluster" "aegis_cluster" {
  name     = var.gke_cluster_name
  location = var.region
  network  = google_compute_network.aegis_vpc.id
  subnetwork = google_compute_subnetwork.aegis_subnet.id

  enable_autopilot = true
  release_channel {
    channel = "REGULAR"
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  resource_labels = {
    project     = "aegis"
    environment = var.environment
  }
}

resource "google_sql_database_instance" "aegis_db_instance" {
  name             = "aegis-postgres-${var.environment}"
  database_version = "POSTGRES_16"
  region           = var.region

  settings {
    tier = var.db_tier
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.aegis_vpc.id
    }
  }
  deletion_protection = true
}

resource "google_sql_database" "aegis_db" {
  name     = "aegis"
  instance = google_sql_database_instance.aegis_db_instance.name
}

resource "google_sql_user" "aegis_db_user" {
  name     = "aegis_user"
  instance = google_sql_database_instance.aegis_db_instance.name
  password = "PLACEHOLDER_CHANGE_ME" # Use secret manager in production
}

resource "google_redis_instance" "aegis_redis" {
  name           = "aegis-redis-${var.environment}"
  tier           = "STANDARD_HA"
  memory_size_gb = var.redis_memory_size_gb
  region         = var.region

  authorized_network = google_compute_network.aegis_vpc.id
  auth_enabled       = true
  transit_encryption_mode = "SERVER_AUTHENTICATION"

  labels = {
    project     = "aegis"
    environment = var.environment
  }
}

locals {
  topics = [
    "aegis-events-raw",
    "aegis-events-detections",
    "aegis-events-rag",
    "aegis-control-policy-reload",
    "aegis-redteam-jobs"
  ]
}

resource "google_pubsub_topic" "topics" {
  for_each = toset(local.topics)
  name     = "${each.key}-${var.environment}"
  labels = {
    project     = "aegis"
    environment = var.environment
  }
}

resource "google_pubsub_subscription" "subscriptions" {
  for_each = toset(local.topics)
  name     = "${each.key}-sub-${var.environment}"
  topic    = google_pubsub_topic.topics[each.key].name

  ack_deadline_seconds = 60
  message_retention_duration = "604800s" # 7 days
}

resource "google_artifact_registry_repository" "aegis_registry" {
  location      = var.region
  repository_id = "aegis-registry-${var.environment}"
  description   = "Docker repository for AEGIS containers"
  format        = "DOCKER"
  labels = {
    project     = "aegis"
    environment = var.environment
  }
}

resource "google_secret_manager_secret" "secrets" {
  for_each = toset(["db-password", "redis-auth-token", "kafka-password"])
  secret_id = "aegis-${each.key}-${var.environment}"
  replication {
    auto {}
  }
  labels = {
    project     = "aegis"
    environment = var.environment
  }
}
