data "archive_file" "cusoon_results_processor_function_distribution" {
  source_file = "../functions/cusoon-results-processor/main"
  output_path = "../functions/cusoon-results-processor/${var.product}-results-processor.zip"
  type        = "zip"
}

data "aws_s3_bucket" "cusoon_suplier_results_bucket" {
  bucket = "${var.product}-${var.environment}-supplier-results-bucket"
}

data "aws_iam_policy_document" "cusoon_results_processor_policy_document" {
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:Query"
    ]
    resources = [
      var.cusoon_results_table_arn,
      "${var.cusoon_results_table_arn}/index/*",
      var.cusoon_email_sends_arn,
      "${var.cusoon_email_sends_arn}/index/*"
    ]
  }

  statement {
    effect    = "Allow"
    actions   = [
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueUrl",
      "sqs:GetQueueAttributes"
    ]
    resources = [
      var.cusoon_supplier_results_queue_arn
    ]
  }

  statement {
    effect    = "Allow"
    actions   = [
      "s3:GetObject"
    ]
    resources = [
      "${data.aws_s3_bucket.cusoon_suplier_results_bucket.arn}/*"
    ]
  }

  statement {
    effect    = "Allow"
    actions   = [
      "s3:PutObject"
    ]
    resources = [
      "${var.cusoon_results_datalake_bucket_arn}/*"
    ]
  }


  statement {
    effect    = "Allow"
    actions   = [
      "ses:SendEmail"
    ]
    resources = [
      "arn:aws:ses:${var.region}:${var.account_number}:identity/${var.ses_domain}"
    ]
  }
}

data "aws_ssm_parameter" "new_relic_license_key" {
  name = "/ssi/${var.environment}/newrelic/license/key"
}