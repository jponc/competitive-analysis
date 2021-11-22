locals {
  service_name_env       = "${var.service_name}-${var.environment}"
  s3_results_bucket_name = "${var.service_name}-${var.environment}-results"
}
