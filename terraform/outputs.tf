output "api_endpoint_url" {
  value = module.lambda.api_endpoint_url
}

output "sqs_queue_name" {
  value = module.sqs.attendee_input_queue_name
}