output "clams_api_endpoint_url" {
  value = module.lambda.clams_api_endpoint_url
}

output "sqs_queue_name" {
  value = module.sqs.attendee_input_queue_name
}

output "clams_frontend_bucket_id" {
  value = module.s3.public_website_bucket_id
}