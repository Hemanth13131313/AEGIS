terraform {
  required_version = ">= 1.8"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

locals {
  tags = {
    project    = "aegis"
    environment = var.environment
    managed_by  = "terraform"
  }
}

# Placeholder VPC
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.0.0"

  name = "aegis-vpc-${var.environment}"
  cidr = var.vpc_cidr

  azs             = ["${var.region}a", "${var.region}b"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]

  enable_nat_gateway = true
  single_nat_gateway = true

  tags = local.tags
}

# PostgreSQL Database
resource "aws_db_instance" "postgres" {
  identifier           = "aegis-db-${var.environment}"
  instance_class       = "db.t3.medium"
  allocated_storage    = 20
  engine               = "postgres"
  engine_version       = "16"
  username             = "aegis_admin"
  password             = var.db_password
  db_subnet_group_name = aws_db_subnet_group.default.name
  vpc_security_group_ids = [aws_security_group.db.id]
  skip_final_snapshot  = true
  multi_az             = false # For dev

  tags = local.tags
}

resource "aws_db_subnet_group" "default" {
  name       = "aegis-db-subnet-group-${var.environment}"
  subnet_ids = module.vpc.private_subnets
  tags       = local.tags
}

resource "aws_security_group" "db" {
  name_prefix = "aegis-db-sg-"
  vpc_id      = module.vpc.vpc_id
  
  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  tags = local.tags
}

# Redis Cluster
resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "aegis-redis-${var.environment}"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  engine_version       = "7.1"
  port                 = 6379
  security_group_ids   = [aws_security_group.redis.id]
  subnet_group_name    = aws_elasticache_subnet_group.default.name
  
  tags = local.tags
}

resource "aws_elasticache_subnet_group" "default" {
  name       = "aegis-redis-subnet-group-${var.environment}"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_security_group" "redis" {
  name_prefix = "aegis-redis-sg-"
  vpc_id      = module.vpc.vpc_id
  
  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  tags = local.tags
}

# Kafka (MSK)
resource "aws_msk_cluster" "kafka" {
  cluster_name           = "aegis-kafka-${var.environment}"
  kafka_version          = "3.5.1"
  number_of_broker_nodes = 2

  broker_node_group_info {
    instance_type   = "kafka.t3.small"
    client_subnets  = module.vpc.private_subnets
    security_groups = [aws_security_group.kafka.id]
  }

  tags = local.tags
}

resource "aws_security_group" "kafka" {
  name_prefix = "aegis-kafka-sg-"
  vpc_id      = module.vpc.vpc_id
  
  ingress {
    from_port   = 9092
    to_port     = 9092
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  tags = local.tags
}

# EKS Cluster stub
resource "aws_eks_cluster" "main" {
  name     = "aegis-eks-${var.environment}"
  version  = "1.30"
  role_arn = aws_iam_role.eks_cluster_role.arn

  vpc_config {
    subnet_ids = module.vpc.private_subnets
  }

  tags = local.tags
}

resource "aws_iam_role" "eks_cluster_role" {
  name = "aegis-eks-cluster-role-${var.environment}"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "eks.amazonaws.com"
        }
      }
    ]
  })
}

# TODO Phase 7 — GCP and Azure modules mirror this structure
