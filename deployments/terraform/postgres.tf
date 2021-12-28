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
    min_capacity             = 4
    seconds_until_auto_pause = 3600 // Shut down if there's no activity after 1hour
    timeout_action           = "ForceApplyCapacityChange"
  }

  db_subnet_group_name   = aws_db_subnet_group.postgres_subnet_group.name
  vpc_security_group_ids = [data.aws_ssm_parameter.vpc_default_security_group_ssm.value]
}

# Secrets manager -- This enables secrets manager to store db credentials which will later be used to access the database through Data API
resource "aws_secretsmanager_secret" "db_instance_credentials_secret" {
  name        = "${var.service_name}-${var.environment}-secret"
  description = ""
}

resource "aws_secretsmanager_secret_version" "db_instance_credentials_secret_values" {
  secret_id     = aws_secretsmanager_secret.db_instance_credentials_secret.id
  secret_string = jsonencode({
    username: aws_rds_cluster.postgres-sls-cluster.master_username,
    password: random_string.postgres_password.result,
    engine: "postgres",
    host: aws_rds_cluster.postgres-sls-cluster.endpoint,
    port: 5432,
    dbInstanceIdentifier: aws_rds_cluster.postgres-sls-cluster.id
  })
}
