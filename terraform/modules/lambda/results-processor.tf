resource "aws_iam_role" "cusoon_results_processor_iam_role" {
  name                  = "${var.product}-${var.environment}-results-processor-iam-role"
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

resource "aws_s3_bucket_object" "cusoon_results_processor_lambda_function_distribution_bucket_object" {
  bucket = "${var.environment}-${var.distribution_bucket}"
  key    = "${var.product}-results-processor}/${var.product}-results-processor.zip"
  source = data.archive_file.cusoon_results_processor_function_distribution.output_path
  etag = filemd5(data.archive_file.cusoon_results_processor_function_distribution.output_path)
}

resource "aws_lambda_function" "cusoon_results_processor_lambda_function" {
  function_name    = "${var.product}-${var.environment}-results-processor"
  handler          = "main"
  runtime          = "go1.x"
  role             = aws_iam_role.cusoon_results_processor_iam_role.arn
  s3_bucket = aws_s3_bucket_object.cusoon_results_processor_lambda_function_distribution_bucket_object.bucket
  s3_key = aws_s3_bucket_object.cusoon_results_processor_lambda_function_distribution_bucket_object.key
  source_code_hash = data.archive_file.cusoon_results_processor_function_distribution.output_md5

  environment {
    variables = {
      FRONTEND_URL      = var.frontend_url
      FROM_EMAIL        = "no-reply@${var.ses_domain}"
      SEND_EMAILS       = var.send_emails
      TO_EMAIL_OVERRIDE = var.to_email_override
      NEW_RELIC_LICENSE_KEY = data.aws_ssm_parameter.new_relic_license_key.value
      ENABLED_ERROR_CODES = var.enabled_error_codes
      RESULTS_DATALAKE_BUCKET_NAME = "${var.product}-${var.environment}-results-datalake"
    }
  }

  tags = {
    Name          = "${var.product}.${var.environment}.lambda.cusoon_results_processor"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    SubProduct    = var.sub_product
    CostCode      = var.cost_code
    Orchestration = var.orchestration
    Description   = "CU Soon Lambda for processing supplier results"
  }
}


resource "aws_iam_role_policy_attachment" "cusoon_results_processor_policy_attachment_execution" {
  role       = aws_iam_role.cusoon_results_processor_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "cusoon_results_processor_iam_policy" {
  name   = "${var.product}-${var.environment}-results-processor-iam-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.cusoon_results_processor_policy_document.json
}

resource "aws_iam_role_policy_attachment" "cusoon_results_processor_policy_attachment" {
  role       = aws_iam_role.cusoon_results_processor_iam_role.name
  policy_arn = aws_iam_policy.cusoon_results_processor_iam_policy.arn
}

resource aws_cloudwatch_log_group "cusoon_results_processor_log_group" {
  name = "/aws/lambda/${aws_lambda_function.cusoon_results_processor_lambda_function.function_name}"
  retention_in_days = 14
}

resource "aws_lambda_event_source_mapping" "cusoon_results_processor_source_mapping" {
  event_source_arn = var.cusoon_supplier_results_queue_arn
  function_name = aws_lambda_function.cusoon_results_processor_lambda_function.arn
}