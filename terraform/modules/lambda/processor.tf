resource "aws_iam_role" "processor_iam_role" {
  name                  = "${var.product}-${var.environment}-processor-iam-role"
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

resource "aws_cloudwatch_log_group" "processor_cloudwatch_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.processor_lambda_function.function_name}"
  retention_in_days = 1
}

data "aws_iam_policy_document" "processor_iam_policy_document" {
  statement {
    effect  = "Allow"
    actions = [
      "dynamodb:PutItem",
    ]
    resources = [
      var.attendees_table_arn
    ]
  }
  statement {
    effect  = "Allow"
    actions = [
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueUrl",
      "sqs:GetQueueAttributes",
    ]
    resources = [
      var.signups_queue_arn
    ]
  }
}

resource "aws_iam_role_policy_attachment" "processor_policy_attachment" {
  role       = aws_iam_role.processor_iam_role.name
  policy_arn = aws_iam_policy.processor_iam_policy.arn
}

resource "aws_iam_role_policy_attachment" "processor_policy_attachment_execution" {
  role       = aws_iam_role.processor_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "processor_iam_policy" {
  name   = "${var.product}-${var.environment}-processor-iam-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.processor_iam_policy_document.json
}

data "archive_file" "processor_lambda_function_distribution" {
  source_file = "../lambdas/processor/main"
  output_path = "../lambdas/processor/${var.product}-processor.zip"
  type        = "zip"
}

resource "aws_s3_object" "processor_lambda_function_distribution_bucket_object" {
  bucket = var.distribution_bucket
  key    = "lambdas/${var.product}-processor/${var.product}-processor.zip"
  source = data.archive_file.processor_lambda_function_distribution.output_path
  etag   = filemd5(data.archive_file.processor_lambda_function_distribution.output_path)
}

resource "aws_lambda_function" "processor_lambda_function" {
  function_name    = "${var.product}-${var.environment}-processor"
  role             = aws_iam_role.processor_iam_role.arn
  handler          = "main"
  runtime          = "go1.x"
  s3_bucket        = aws_s3_object.processor_lambda_function_distribution_bucket_object.bucket
  s3_key           = aws_s3_object.processor_lambda_function_distribution_bucket_object.key
  source_code_hash = data.archive_file.processor_lambda_function_distribution.output_md5
  timeout          = 60
  memory_size      = 256

  environment {
    variables = {
      ATTENDEES_TABLE_NAME = var.attendees_table_name
      AWS_REGION           = var.region
      DB_HOST              = var.db_host,
      DB_NAME              = var.db_name,
      DB_USER              = var.db_username,
      DB_PASSWORD          = var.db_password,
    }
  }

  tags = {
    Name          = "${var.product}.${var.environment}.lambda.processor"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    Orchestration = var.orchestration
    Description   = "Lambda for processing sign-up requests incoming to the database"
  }
}

resource "aws_lambda_event_source_mapping" "processor_event_source_mapping" {
  event_source_arn        = var.signups_queue_arn
  function_name           = aws_lambda_function.processor_lambda_function.arn
  function_response_types = ["ReportBatchItemFailures"]
}


