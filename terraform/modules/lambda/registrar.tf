resource "aws_iam_role" "registrar_iam_role" {
  name                  = "${var.product}-${var.environment}-registrar-iam-role"
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

data "archive_file" "registrar_lambda_function_distribution" {
  source_dir  = "../functions/registrar"
  output_path = "../functions/registrar/${var.product}-registrar.zip"
  type        = "zip"
}

resource "aws_s3_bucket_object" "registrar_lambda_function_distribution_bucket_object" {
  bucket = "${var.environment}-${var.distribution_bucket}"
  key    = "${var.product}-registrar}/${var.product}-registrar.zip"
  source = data.archive_file.registrar_lambda_function_distribution.output_path
  etag   = filemd5(data.archive_file.registrar_lambda_function_distribution.output_path)
}

resource "aws_lambda_function" "registrar_lambda_function" {
  function_name    = "${var.product}-${var.environment}-registrar"
  handler          = "main"
  runtime          = "go1.x"
  role             = aws_iam_role.registrar_iam_role.arn
  s3_bucket        = aws_s3_bucket_object.registrar_lambda_function_distribution_bucket_object.bucket
  s3_key           = aws_s3_bucket_object.registrar_lambda_function_distribution_bucket_object.key
  source_code_hash = data.archive_file.registrar_lambda_function_distribution.output_md5

  environment {
    variables = {
      ATTENDEES_TABLE_NAME = var.attendees_table_name
      INPUT_QUEUE_NAME     = var.input_queue_name
    }
  }

  tags = {
    Name          = "${var.product}.${var.environment}.lambda.registrar"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    SubProduct    = var.sub_product
    CostCode      = var.cost_code
    Orchestration = var.orchestration
    Description   = "registrar for processing incoming eHAMS attendees"
  }
}

data "aws_iam_policy_document" "registrar_iam_policy_document" {
  statement {
    effect  = "Allow"
    actions = [
      "ec2:DescribeNetworkInterfaces",
      "dynamodb:UpdateItem",
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
      var.input_queue_arn
    ]
  }
}

resource "aws_iam_role_policy_attachment" "registrar_policy_attachment" {
  role       = aws_iam_role.registrar_iam_role.name
  policy_arn = aws_iam_policy.registrar_iam_policy.arn
}

resource "aws_iam_role_policy_attachment" "registrar_policy_attachment_execution" {
  role       = aws_iam_role.registrar_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "registrar_iam_policy" {
  name   = "${var.environment}-ope-test-mbh-ehams-registrar-iam-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.registrar_iam_policy_document.json
}

resource aws_cloudwatch_log_group "registrar_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.registrar_lambda_function.function_name}"
  retention_in_days = 1
}

resource "aws_lambda_event_source_mapping" "registrar_source_mapping" {
  event_source_arn        = var.input_queue_arn
  function_name           = aws_lambda_function.registrar_lambda_function.arn
  function_response_types = ["ReportBatchItemFailures"]
}