output "clams_api_endpoint_url" {
  value = module.lambda.clams_api_endpoint_url
}

output "sqs_queue_name" {
  value = module.sqs.attendee_input_queue_name
}

output "clams_frontend_bucket_id" {
  value = module.s3.public_website_bucket_id
}

output "rds_database_host" {
  value = module.rds.rds_database_host
}

output "rds_database_name" {
  value = module.rds.rds_database_name
}

output "rds_database_username" {
  value     = module.rds.rds_database_username
  sensitive = true
}

output "rds_database_password" {
  value     = module.rds.rds_database_password
  sensitive = true
}