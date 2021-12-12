variable "environment" {
  default     = "staging"
  description = "Environment"
}

variable "service_name" {
  description = "This is the name of the service"
}

variable "aws_profile" {
  description = "AWS Profile"
}

variable "aws_region" {
  description = "AWS Region"
}

variable "rds_database_name" {
  type        = string
  description = "RDS database name"
}

variable "private_subnets_ssm_path" {
  type        = string
  description = "SSM Path to Private subnets StringList"
}

variable "domain_ssm_value" {
  type        = string
  description = "SSM Domain Value"
}

variable "vpc_default_security_group_ssm_path" {
  type        = string
  description = "SSM Path to VPC default security group"
}

variable "appextend_certificate_arn_ssm_path" {
  type        = string
  description = "Appextend Certificate ARN Path"
}

variable "appextend_hosted_zone_id_ssm_path" {
  type        = string
  description = "Appextend Route 53 Hosted Zone ID Path"
}

variable "frontend_url" {
  type        = string
  description = "Frontend URL"
}
