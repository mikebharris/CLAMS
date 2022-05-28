resource "aws_iam_role" "cusoon_results_api_iam_role" {
  name                  = "${var.product}-${var.environment}-results-api-iam-role"
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

resource "aws_cloudwatch_log_group" "cusoon_results_api_cloudwatch_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.cusoon_results_api_lambda_function.function_name}"
  retention_in_days = 14
}

data "aws_iam_policy_document" "cusoon_results_api_iam_policy_document" {
  statement {
    effect  = "Allow"
    actions = [
      "ec2:DescribeNetworkInterfaces",
      "dynamodb:GetItem",
      "dynamodb:Query"
    ]
    resources = [
      "arn:aws:dynamodb:${var.region}:${var.account_number}:table/cusoon-results-datastore"
    ]
  }
  statement {
    effect  = "Allow"
    actions = [
      "dynamodb:UpdateItem",
    ]
    resources = [
      "arn:aws:dynamodb:${var.region}:${var.account_number}:table/cusoon-email-sends"
    ]
  }
  statement {
    effect  = "Allow"
    actions = [
      "ec2:DescribeNetworkInterfaces",
      "ec2:CreateNetworkInterface",
      "ec2:DeleteNetworkInterface"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy_attachment" "cusoon_results_api_policy_attachment" {
  role       = aws_iam_role.cusoon_results_api_iam_role.name
  policy_arn = aws_iam_policy.cusoon_results_api_iam_policy.arn
}

resource "aws_iam_role_policy_attachment" "cusoon_results_api_policy_attachment_execution" {
  role       = aws_iam_role.cusoon_results_api_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "cusoon_results_api_iam_policy" {
  name   = "${var.environment}-cusoon-results-api-iam-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.cusoon_results_api_iam_policy_document.json
}

data "archive_file" "cusoon_results_api_lambda_function_distribution" {
  source_dir  = "../functions/cusoon-results-api"
  output_path = "../functions/cusoon-results-api/${var.product}-results-api.zip"
  type        = "zip"
}

data "aws_s3_bucket" "cusoon_results_api_lambda_function_distribution_bucket" {
  bucket = "${var.environment}-${var.distribution_bucket}"
}

resource "aws_s3_bucket_object" "cusoon_results_api_lambda_function_distribution_bucket_object" {
  bucket = "${var.environment}-${var.distribution_bucket}"
  key    = "${var.product}-results-api/${var.product}-results-api.zip"
  source = data.archive_file.cusoon_results_api_lambda_function_distribution.output_path
  etag   = filemd5(data.archive_file.cusoon_results_api_lambda_function_distribution.output_path)
}

resource "aws_lambda_function" "cusoon_results_api_lambda_function" {
  function_name    = "${var.product}-${var.environment}-results-api"
  role             = aws_iam_role.cusoon_results_api_iam_role.arn
  handler          = "main"
  runtime          = "go1.x"
  s3_bucket        = data.aws_s3_bucket.cusoon_results_api_lambda_function_distribution_bucket.bucket
  s3_key           = aws_s3_bucket_object.cusoon_results_api_lambda_function_distribution_bucket_object.key
  source_code_hash = data.archive_file.cusoon_results_api_lambda_function_distribution.output_md5
  timeout          = 60
  memory_size      = 256

  tags = {
    Name          = "${var.product}.${var.environment}.lambda.cusoon_results_processor"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    SubProduct    = var.sub_product
    CostCode      = var.cost_code
    Orchestration = var.orchestration
    Description   = "Lambda for fetching the CU Soon results for a given unique Submission UUID and passing it to the front-end UI service"
  }
}


