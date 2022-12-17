resource "aws_iam_role" "attendees_api_iam_role" {
  name                  = "${var.product}-${var.environment}-attendees-api-iam-role"
  force_detach_policies = true
  assume_role_policy    = jsonencode({
    Version   = "2012-10-17"
    Statement = [
      {
        Action    = "sts:AssumeRole"
        Effect    = "Allow"
        Sid       = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_cloudwatch_log_group" "attendees_api_cloudwatch_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.attendees_api_lambda_function.function_name}"
  retention_in_days = 1
}

data "aws_iam_policy_document" "attendees_api_iam_policy_document" {
  statement {
    effect  = "Allow"
    actions = [
      "ec2:DescribeNetworkInterfaces",
      "dynamodb:GetItem",
      "dynamodb:Scan"
    ]
    resources = [
      var.attendees_table_arn
    ]
  }
}

resource "aws_iam_role_policy_attachment" "attendees_api_policy_attachment" {
  role       = aws_iam_role.attendees_api_iam_role.name
  policy_arn = aws_iam_policy.attendees_api_iam_policy.arn
}

resource "aws_iam_role_policy_attachment" "attendees_api_policy_attachment_execution" {
  role       = aws_iam_role.attendees_api_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "attendees_api_iam_policy" {
  name   = "${var.product}-${var.environment}-attendees-api-iam-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.attendees_api_iam_policy_document.json
}

data "archive_file" "attendees_api_lambda_function_distribution" {
  source_file = "../functions/attendees-api/main"
  output_path = "../functions/attendees-api/${var.product}-attendees-api.zip"
  type        = "zip"
}

resource "aws_s3_object" "attendees_api_lambda_function_distribution_bucket_object" {
  bucket = "${var.account_number}-${var.distribution_bucket}"
  key    = "lambdas/${var.product}-attendees-api/${var.product}-attendees-api.zip"
  source = data.archive_file.attendees_api_lambda_function_distribution.output_path
  etag   = filemd5(data.archive_file.attendees_api_lambda_function_distribution.output_path)
}

resource "aws_lambda_function" "attendees_api_lambda_function" {
  function_name    = "${var.product}-${var.environment}-attendees-api"
  role             = aws_iam_role.attendees_api_iam_role.arn
  handler          = "main"
  runtime          = "go1.x"
  s3_bucket        = aws_s3_object.attendees_api_lambda_function_distribution_bucket_object.bucket
  s3_key           = aws_s3_object.attendees_api_lambda_function_distribution_bucket_object.key
  source_code_hash = data.archive_file.attendees_api_lambda_function_distribution.output_md5
  timeout          = 60
  memory_size      = 256

  environment {
    variables = {
      ATTENDEES_TABLE_NAME = var.attendees_table_name
    }
  }

  tags = {
    Name          = "${var.product}.${var.environment}.lambda.attendees-api"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    Orchestration = var.orchestration
    Description   = "Lambda for fetching attendees from DynamoDB and displaying them"
  }
}


