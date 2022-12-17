#resource "aws_iam_role" "api_authorizer_iam_role" {
#  name                  = "${var.product}-${var.environment}-api-authorizer-iam-role"
#  force_detach_policies = true
#  assume_role_policy    = jsonencode({
#    Version   = "2012-10-17"
#    Statement = [
#      {
#        Action    = "sts:AssumeRole"
#        Effect    = "Allow"
#        Sid       = ""
#        Principal = {
#          Service = "lambda.amazonaws.com"
#        }
#      }
#    ]
#  })
#}
#
#resource aws_cloudwatch_log_group "api_authorizer_log_group" {
#  name              = "/aws/lambda/${aws_lambda_function.api_authorizer_lambda_function.function_name}"
#  retention_in_days = 14
#}
#
#resource "aws_iam_role_policy_attachment" "api_authorizer_policy_attachment" {
#  role       = aws_iam_role.api_authorizer_iam_role.name
#  policy_arn = aws_iam_policy.api_authorizer_listener_iam_policy.arn
#}
#
#resource "aws_iam_role_policy_attachment" "api_authorizer_policy_attachment_execution" {
#  role       = aws_iam_role.api_authorizer_iam_role.name
#  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
#}
#
#data "aws_iam_policy_document" "authorizer_lambda_access_policy_document" {
#  statement {
#    effect  = "Allow"
#    actions = [
#      "ssm:GetParameter"
#    ]
#    resources = [
#      "arn:aws:ssm:${var.region}:${var.account_number}:parameter/${var.product}/${var.environment}/api/user",
#      "arn:aws:ssm:${var.region}:${var.account_number}:parameter/${var.product}/${var.environment}/api/password"
#    ]
#  }
#
#  statement {
#    effect = "Allow"
#
#    actions = [
#      "logs:CreateLogGroup",
#      "ec2:CreateNetworkInterface",
#      "ec2:DescribeNetworkInterfaces",
#      "ec2:DeleteNetworkInterface"
#    ]
#
#    resources = [
#      "*"
#    ]
#  }
#}
#
#resource "aws_iam_policy" "api_authorizer_listener_iam_policy" {
#  name   = "${var.product}-${var.environment}-api-authorizer-iam-policy"
#  path   = "/"
#  policy = data.aws_iam_policy_document.authorizer_lambda_access_policy_document.json
#}
#
#data "archive_file" "api_authorizer_lambda_function_distribution" {
#  source_file = "../functions/authorizer/main"
#  output_path = "../functions/authorizer/authorizer.zip"
#  type        = "zip"
#}
#
#resource "aws_s3_object" "api_authorizer_lambda_function_distribution_s3_object" {
#  bucket = "${var.environment}-${var.distribution_bucket}"
#  key    = "${var.product}-api/api-authorizer.zip"
#  source = data.archive_file.api_authorizer_lambda_function_distribution.output_path
#  etag   = filemd5(data.archive_file.api_authorizer_lambda_function_distribution.output_path)
#}
#
#resource "aws_lambda_function" "api_authorizer_lambda_function" {
#  function_name    = "${var.product}-${var.environment}-api-authorizer"
#  handler          = "main"
#  runtime          = "go1.x"
#  memory_size      = 128
#  role             = aws_iam_role.api_authorizer_iam_role.arn
#  s3_bucket        = aws_s3_object.api_authorizer_lambda_function_distribution_s3_object.bucket
#  s3_key           = aws_s3_object.api_authorizer_lambda_function_distribution_s3_object.key
#  source_code_hash = data.archive_file.api_authorizer_lambda_function_distribution.output_md5
#
#  environment {
#    variables = {
#      ENVIRONMENT = var.environment
#    }
#  }
#
#  tags = {
#    Name          = "${var.product}.${var.environment}.lambda.api_authorizer"
#    Contact       = var.contact
#    Environment   = var.environment
#    Product       = var.product
#    Orchestration = var.orchestration
#    Description   = "Lambda for authorising requests to the CLAMS API"
#  }
#}
#
#resource "aws_apigatewayv2_authorizer" "api_authorizer_authorizer" {
#  api_id          = aws_apigatewayv2_api.clams_api_lambda_http_gateway.id
#  authorizer_type = "REQUEST"
#  name            = "${var.product}-${var.environment}-api-authorizer"
#  authorizer_uri = aws_lambda_function.api_authorizer_lambda_function.invoke_arn
#  identity_sources = ["$request.header.Authorization"]
#  authorizer_payload_format_version = "2.0"
#  enable_simple_responses = true
#}
#
#resource "aws_lambda_permission" "api_authorizer_http_gateway_lambda_permission" {
#  action        = "lambda:InvokeFunction"
#  function_name = aws_lambda_function.api_authorizer_lambda_function.function_name
#  principal     = "apigateway.amazonaws.com"
#  statement_id  = "AllowExecutionFromAPIGateway"
#  source_arn = "${aws_apigatewayv2_api.clams_api_lambda_http_gateway.execution_arn}/*/*"
#}
#
#resource "aws_apigatewayv2_route" "api_lambda_http_gateway_integration_put_subscription_method_route" {
#  api_id             = aws_apigatewayv2_api.clams_api_lambda_http_gateway.id
#  route_key          = "GET /journal/{journal}/document/{document}/metadata"
#  target             = "integrations/${aws_apigatewayv2_integration.attendees_api_lambda_http_gateway_integration.id}"
#  authorizer_id      = aws_apigatewayv2_authorizer.api_authorizer_authorizer.id
#  authorization_type = "CUSTOM"
#}
#
#
