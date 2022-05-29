resource "aws_apigatewayv2_api" "attendees_api_lambda_http_gateway" {
  name          = "${var.product}-${var.environment}-attendees-api-http-gateway"
  description   = "HTTP API Gateway for MBH Test eHAMS API Lambda."
  protocol_type = "HTTP"

  dynamic "cors_configuration" {
    for_each = length(keys(var.cors_configuration)) == 0 ? [] : [var.cors_configuration]

    content {
      allow_credentials = lookup(cors_configuration.value, "allow_credentials", null)
      allow_headers     = lookup(cors_configuration.value, "allow_headers", null)
      allow_methods     = lookup(cors_configuration.value, "allow_methods", null)
      allow_origins     = lookup(cors_configuration.value, "allow_origins", null)
      expose_headers    = lookup(cors_configuration.value, "expose_headers", null)
      max_age           = lookup(cors_configuration.value, "max_age", null)
    }
  }
}

resource "aws_apigatewayv2_stage" "attendees_api_lambda_http_gateway_stage" {
  api_id      = aws_apigatewayv2_api.attendees_api_lambda_http_gateway.id
  name        = "attendees"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.attendees_api_http_gateway_log_group.arn

    format          = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
    })
  }
}

resource "aws_apigatewayv2_integration" "attendees_api_lambda_http_gateway_integration_get_method" {
  api_id           = aws_apigatewayv2_api.attendees_api_lambda_http_gateway.id
  integration_type = "AWS_PROXY"
  integration_method = "POST"
  integration_uri = aws_lambda_function.attendees_api_lambda_function.invoke_arn
}

resource "aws_apigatewayv2_route" "attendees_api_lambda_http_gateway_integration_get_method_route" {
  api_id    = aws_apigatewayv2_api.attendees_api_lambda_http_gateway.id
  route_key = "GET /{code}"
  target = "integrations/${aws_apigatewayv2_integration.attendees_api_lambda_http_gateway_integration_get_method.id}"
}

resource "aws_cloudwatch_log_group" "attendees_api_http_gateway_log_group" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.attendees_api_lambda_http_gateway.name}"
  retention_in_days = 1
}

resource "aws_lambda_permission" "attendees_api_http_gateway_lambda_permission" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.attendees_api_lambda_function.function_name
  principal     = "apigateway.amazonaws.com"
  statement_id = "AllowExecutionFromAPIGateway"

  source_arn = "${aws_apigatewayv2_api.attendees_api_lambda_http_gateway.execution_arn}/*/*"
}

resource "aws_ssm_parameter" "attendees_api_endpoint" {
  name  = "/${var.product}/${var.environment}/attendees-api/endpoint"
  type  = "String"
  value = aws_apigatewayv2_api.attendees_api_lambda_http_gateway.api_endpoint
  overwrite = true
}