output "cloudfront_domain_name" {
  value = aws_cloudfront_distribution.public_site_cdn.domain_name
}

output "cloudfront_distribution_id" {
  value = aws_cloudfront_distribution.public_site_cdn.id
}

output "origin_access_identity" {
  value = aws_cloudfront_origin_access_identity.origin_access_identity
}
