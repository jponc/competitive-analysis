resource "aws_s3_bucket" "frontend_bucket" {
  bucket = var.frontend_url
  acl    = "public-read"
  policy = <<-EOT
  {
      "Version": "2012-10-17",
      "Statement": [
          {
              "Sid": "PublicRead",
              "Effect": "Allow",
              "Principal": "*",
              "Action": "s3:GetObject",
              "Resource": "arn:aws:s3:::${var.frontend_url}/*"
          }
      ]
  }
  EOT

  website {
    index_document = "index.html"
    error_document = "index.html"
  }
}

locals {
  s3_origin_id = "S3-${aws_s3_bucket.frontend_bucket.id}"
}

resource "aws_cloudfront_distribution" "frontend_distribution" {
  origin {
    domain_name = aws_s3_bucket.frontend_bucket.bucket_regional_domain_name
    origin_id   = local.s3_origin_id
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  aliases = [var.frontend_url]

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  # Default behaviour
  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id

    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
    viewer_protocol_policy = "allow-all"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  # Cache behavior with precedence 0
  ordered_cache_behavior {
    path_pattern     = "*"
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id

    min_ttl                = 0
    default_ttl            = 86400
    max_ttl                = 31536000
    viewer_protocol_policy = "redirect-to-https"
    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  # Custom error response to handle 4xx
  custom_error_response {
    error_code = 403
    error_caching_min_ttl = 60
    response_code = 200
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code = 404
    error_caching_min_ttl = 60
    response_code = 200
    response_page_path = "/index.html"
  }


  # Certificate
  viewer_certificate {
    acm_certificate_arn      = data.aws_ssm_parameter.appextend_certificate_arn_ssm.value
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2019"
  }
}

resource "aws_route53_record" "frontend_route_53" {
  zone_id = data.aws_ssm_parameter.appextend_hosted_zone_id_ssm.value
  name    = var.frontend_url
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.frontend_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.frontend_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}
