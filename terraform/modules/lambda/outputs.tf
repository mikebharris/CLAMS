output "clams_api_endpoint_url" {
  value = aws_apigatewayv2_api.clams_api_lambda_http_gateway.api_endpoint
}

output "clams_api_lambda_function_url" {
  value = aws_lambda_function_url.attendees_api_lambda_function_url.function_url
}