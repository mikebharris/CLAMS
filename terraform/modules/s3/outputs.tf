output "public_website_bucket_id" {
  value = aws_s3_bucket.clams_frontend_static_content_bucket.id
}

output "bucket_regional_domain_name" {
  value = aws_s3_bucket.clams_frontend_static_content_bucket.bucket_regional_domain_name
}