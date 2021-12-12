environment  = "production"
service_name = "competitive-analysis"
aws_profile  = "uptactics"
aws_region   = "us-east-1"

# RDS
rds_database_name     = "mydatabase"

# VPC
private_subnets_ssm_path            = "/uptactics/production/PRIVATE_SUBNET_IDS"
vpc_default_security_group_ssm_path = "/uptactics/production/DEFAULT_SECURITY_GROUP"

# Domain
domain_ssm_value = "api-production-competitive-analysis.appextend.com"

# Frontend
frontend_url                       = "seo-tool.appextend.com"
appextend_hosted_zone_id_ssm_path  = "/uptactics/APPEXTEND_HOSTED_ZONE_ID"
appextend_certificate_arn_ssm_path = "/uptactics/APPEXTEND_CERTIFICATE_ARN"
