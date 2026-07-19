output "vpc_id" {
  value = module.vpc.vpc_id
}

output "db_endpoint" {
  value = aws_db_instance.postgres.endpoint
}

output "redis_endpoint" {
  value = aws_elasticache_cluster.redis.cache_nodes[0].address
}

output "kafka_brokers" {
  value = aws_msk_cluster.kafka.bootstrap_brokers
}

output "eks_endpoint" {
  value = aws_eks_cluster.main.endpoint
}
