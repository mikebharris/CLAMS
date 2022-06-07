resource "aws_cloudfront_origin_access_identity" "origin_access_identity" {
  comment = "${var.environment}-origin-access-identity"
}

resource "aws_cloudfront_distribution" "public_site_cdn" {
  enabled      = true
  price_class  = "PriceClass_200"
  http_version = "http1.1"
  aliases      = [var.frontend_domain]

  origin {
    origin_id   = "S3-${var.frontend_domain}"

    domain_name = var.origin_domain_name

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.origin_access_identity.cloudfront_access_identity_path
    }
  }

  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-${var.frontend_domain}"

    min_ttl     = "0"
    default_ttl = "300"
    max_ttl     = "1200"

    viewer_protocol_policy = "redirect-to-https"
    compress               = true

    forwarded_values {
      query_string = true

      cookies {
        forward = "none"
      }
    }
  }

  custom_error_response {
    error_code         = 404
    response_code      = 404
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 403
    response_code      = 200
    response_page_path = "/index.html"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags = {
    Name          = "${var.product}-${var.environment}.cloudfront_distribution"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    Orchestration = var.orchestration
    Description   = "Cloudfront Distribution fronting S3 bucket"
  }

  viewer_certificate {
    acm_certificate_arn            = var.acm_certificate_arn
    ssl_support_method             = "sni-only"
    minimum_protocol_version       = "TLSv1.2_2018"
    cloudfront_default_certificate = false
  }
}

data "aws_route53_zone" "hacktionlab_r53_zone" {
  name         = var.certificate_domain
  private_zone = false
}
