resource "random_string" "jwt_secret" {
  length           = 16
  special          = true
  override_special = "/@£$"
}

resource "aws_ssm_parameter" "jwt_secret_ssm" {
  name  = "/${var.service_name}/${var.environment}/JWT_SECRET"
  type  = "SecureString"
  value = random_string.jwt_secret.result
}

resource "aws_ssm_parameter" "domain_ssm" {
  name  = "/${var.service_name}/${var.environment}/DOMAIN"
  type  = "String"
  value = var.domain_ssm_value
}

resource "aws_ssm_parameter" "postgres_db_conn_url_ssm" {
  name  = "/${var.service_name}/${var.environment}/DB_CONN_URL"
  type  = "SecureString"
  value = "postgres://${aws_rds_cluster.postgres-sls-cluster.master_username}:${random_string.postgres_password.result}@${aws_rds_cluster.postgres-sls-cluster.endpoint}/${var.rds_database_name}"

  depends_on = [
    aws_rds_cluster.postgres-sls-cluster,
  ]
}
