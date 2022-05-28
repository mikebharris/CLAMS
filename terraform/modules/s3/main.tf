resource "aws_s3_bucket" "cusoon_results_datalake" {
  bucket = "${var.product}-${var.environment}-results-datalake"
  acl    = "private"

  tags = {
    Name          = "${var.product}.${var.environment}.s3.supplier_results"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    SubProduct    = var.sub_product
    CostCode      = var.cost_code
    Orchestration = var.orchestration
    Description   = "S3 bucket for storing results in CSV format for data analysis and reporting purposes."
  }
}

resource "aws_s3_bucket_public_access_block" "cusoon_results_datalake_public_access_block" {
  bucket = aws_s3_bucket.cusoon_results_datalake.id

  block_public_acls        = true
  block_public_policy      = true
  ignore_public_acls       = true
  restrict_public_buckets  = true
}