# Generates a random string for our postgres password
resource "random_string" "postgres_password" {
  length  = 30
  special = false
}

# This is the subnet group of our postgres instance
resource "aws_db_subnet_group" "postgres_subnet_group" {
  name       = "${var.service_name}-${var.environment}-db-subnet-group"
  subnet_ids = split(",", data.aws_ssm_parameter.private_subnets_ssm.value)
}

resource "aws_rds_cluster" "postgres-sls-cluster" {
  cluster_identifier      = "${var.service_name}-${var.environment}-sls"
  engine                  = "aurora-postgresql"
  engine_mode             = "serverless"
  enable_http_endpoint    = true
  database_name           = var.rds_database_name
  master_username         = "postgres"
  master_password         = random_string.postgres_password.result
  backup_retention_period = 5

  scaling_configuration {
    auto_pause               = true
    max_capacity             = 32
    min_capacity             = 2
    seconds_until_auto_pause = 3600
    timeout_action           = "ForceApplyCapacityChange"
  }

  db_subnet_group_name   = aws_db_subnet_group.postgres_subnet_group.name
  vpc_security_group_ids = [data.aws_ssm_parameter.vpc_default_security_group_ssm.value]
}
