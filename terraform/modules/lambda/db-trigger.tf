resource "aws_iam_role" "db_trigger_iam_role" {
  name                  = "${var.product}-${var.environment}-db-trigger-iam-role"
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

resource "aws_cloudwatch_log_group" "db_trigger_cloudwatch_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.db_trigger_lambda_function.function_name}"
  retention_in_days = 1
}

data "aws_iam_policy_document" "db_trigger_iam_policy_document" {
  statement {
    effect  = "Allow"
    actions = [
      "sqs:GetQueueUrl",
      "sqs:SendMessage",
    ]
    resources = [
      var.signups_queue_arn
    ]
  }
}

resource "aws_iam_role_policy_attachment" "db_trigger_policy_attachment" {
  role       = aws_iam_role.db_trigger_iam_role.name
  policy_arn = aws_iam_policy.db_trigger_iam_policy.arn
}

resource "aws_iam_role_policy_attachment" "db_trigger_policy_attachment_execution" {
  role       = aws_iam_role.db_trigger_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "db_trigger_iam_policy" {
  name   = "${var.product}-${var.environment}-db_trigger-iam-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.db_trigger_iam_policy_document.json
}

data "archive_file" "db_trigger_lambda_function_distribution" {
  source_file = "../functions/db-trigger/main"
  output_path = "../functions/db-trigger/${var.product}-db-trigger.zip"
  type        = "zip"
}

resource "aws_s3_object" "db_trigger_lambda_function_distribution_bucket_object" {
  bucket = "${var.account_number}-${var.distribution_bucket}"
  key    = "lambdas/${var.product}-db-trigger/${var.product}-db-trigger.zip"
  source = data.archive_file.db_trigger_lambda_function_distribution.output_path
  etag   = filemd5(data.archive_file.db_trigger_lambda_function_distribution.output_path)
}

resource "aws_lambda_function" "db_trigger_lambda_function" {
  function_name    = "${var.product}-${var.environment}-db-trigger"
  role             = aws_iam_role.db_trigger_iam_role.arn
  handler          = "main"
  runtime          = "go1.x"
  s3_bucket        = aws_s3_object.db_trigger_lambda_function_distribution_bucket_object.bucket
  s3_key           = aws_s3_object.db_trigger_lambda_function_distribution_bucket_object.key
  source_code_hash = data.archive_file.db_trigger_lambda_function_distribution.output_md5
  timeout          = 60
  memory_size      = 256

  environment {
    variables = {
      DB_HOST : var.db_host,
      DB_NAME : var.db_name,
      DB_USER : var.db_username,
      DB_PASSWORD : var.db_password,
    }
  }

  tags = {
    Name          = "${var.product}.${var.environment}.lambda.db-trigger"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    Orchestration = var.orchestration
    Description   = "Lambda for turning database triggers into SQS events"
  }
}

resource "aws_cloudwatch_event_rule" "db_trigger_event_rule" {
  name                = "${var.product}-${var.environment}-db-trigger-frequency-rule"
  description         = "The frequency to run the database trigger Lambda"
  schedule_expression = "cron(* * * * ? *)"
  is_enabled          = false
}

resource "aws_cloudwatch_event_target" "db_trigger_event_target" {
  rule      = aws_cloudwatch_event_rule.db_trigger_event_rule.name
  target_id = "lambda"
  arn       = aws_lambda_function.db_trigger_lambda_function.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_db_trigger_lambda" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.db_trigger_lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.db_trigger_event_rule.arn
}
