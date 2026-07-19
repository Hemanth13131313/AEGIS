provider "azurerm" {
  features {}
  subscription_id = var.subscription_id
}

resource "azurerm_resource_group" "aegis_rg" {
  name     = "aegis-${var.environment}-rg"
  location = var.resource_group_location
  tags = {
    Project     = "aegis"
    Environment = var.environment
  }
}

resource "azurerm_virtual_network" "aegis_vnet" {
  name                = "aegis-vnet-${var.environment}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.aegis_rg.location
  resource_group_name = azurerm_resource_group.aegis_rg.name
}

resource "azurerm_subnet" "aegis_subnet" {
  name                 = "aegis-subnet-${var.environment}"
  resource_group_name  = azurerm_resource_group.aegis_rg.name
  virtual_network_name = azurerm_virtual_network.aegis_vnet.name
  address_prefixes     = ["10.0.0.0/24"]
}

resource "azurerm_kubernetes_cluster" "aegis_aks" {
  name                = var.aks_cluster_name
  location            = azurerm_resource_group.aegis_rg.location
  resource_group_name = azurerm_resource_group.aegis_rg.name
  dns_prefix          = "aegisaks"

  role_based_access_control_enabled = true
  workload_identity_enabled         = true
  oidc_issuer_enabled               = true

  default_node_pool {
    name                = "system"
    vm_size             = "Standard_D4s_v3"
    enable_auto_scaling = true
    min_count           = 2
    max_count           = 10
    vnet_subnet_id      = azurerm_subnet.aegis_subnet.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "azure"
  }

  tags = {
    Project     = "aegis"
    Environment = var.environment
  }
}

resource "azurerm_private_dns_zone" "aegis_dns" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.aegis_rg.name
}

resource "azurerm_postgresql_flexible_server" "aegis_db_server" {
  name                   = "aegis-postgres-${var.environment}"
  resource_group_name    = azurerm_resource_group.aegis_rg.name
  location               = azurerm_resource_group.aegis_rg.location
  version                = "16"
  administrator_login    = "psqladmin"
  administrator_password = "PLACEHOLDER_CHANGE_ME"
  storage_mb             = 32768
  sku_name               = var.db_sku_name
  zone                   = "1"
  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "2"
  }
  tags = {
    Project     = "aegis"
    Environment = var.environment
  }
}

resource "azurerm_postgresql_flexible_server_database" "aegis_db" {
  name      = "aegis"
  server_id = azurerm_postgresql_flexible_server.aegis_db_server.id
  collation = "en_US.utf8"
  charset   = "utf8"
}

resource "azurerm_redis_cache" "aegis_redis" {
  name                = "aegis-redis-${var.environment}"
  location            = azurerm_resource_group.aegis_rg.location
  resource_group_name = azurerm_resource_group.aegis_rg.name
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"
  enable_non_ssl_port = false
  minimum_tls_version = "1.2"

  tags = {
    Project     = "aegis"
    Environment = var.environment
  }
}

resource "azurerm_eventhub_namespace" "aegis_eh_ns" {
  name                = "aegis-ehns-${var.environment}"
  location            = azurerm_resource_group.aegis_rg.location
  resource_group_name = azurerm_resource_group.aegis_rg.name
  sku                 = "Standard"
  capacity            = 1
  tags = {
    Project     = "aegis"
    Environment = var.environment
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

resource "azurerm_eventhub" "topics" {
  for_each            = toset(local.topics)
  name                = each.key
  namespace_name      = azurerm_eventhub_namespace.aegis_eh_ns.name
  resource_group_name = azurerm_resource_group.aegis_rg.name
  partition_count     = 2
  message_retention   = 7
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "aegis_kv" {
  name                        = "aegis-kv-${var.environment}"
  location                    = azurerm_resource_group.aegis_rg.location
  resource_group_name         = azurerm_resource_group.aegis_rg.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  tags = {
    Project     = "aegis"
    Environment = var.environment
  }
}

resource "azurerm_key_vault_secret" "placeholders" {
  for_each     = toset(["db-password", "redis-auth-token", "kafka-password"])
  name         = each.key
  value        = "PLACEHOLDER"
  key_vault_id = azurerm_key_vault.aegis_kv.id
}

resource "azurerm_container_registry" "aegis_acr" {
  name                = "aegisacr${var.environment}"
  resource_group_name = azurerm_resource_group.aegis_rg.location
  location            = azurerm_resource_group.aegis_rg.name
  sku                 = "Standard"
  admin_enabled       = false
  tags = {
    Project     = "aegis"
    Environment = var.environment
  }
}
