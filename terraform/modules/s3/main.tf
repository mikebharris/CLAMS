resource "aws_s3_bucket" "clams_frontend_static_content_bucket" {
  bucket        = var.frontend_domain
  force_destroy = true

  tags = {
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    Orchestration = var.orchestration
    Description   = "S3 bucket for CLAMS frontend static content in ${var.environment} environment."
  }
}

resource "aws_s3_account_public_access_block" "blah" {
  block_public_acls = false
}

resource "aws_s3_bucket_ownership_controls" "clams_frontend_static_content_bucket_ownership_controls" {
  bucket = aws_s3_bucket.clams_frontend_static_content_bucket.id

  rule {
    object_ownership = "ObjectWriter"
  }
}

resource "aws_s3_bucket_policy" "clams_frontend_static_content_bucket_policy" {
  bucket = aws_s3_bucket.clams_frontend_static_content_bucket.id
  policy = data.aws_iam_policy_document.clams_frontend_bucket_policy.json
}

resource "aws_s3_bucket_website_configuration" "clams_frontend_static_content_bucket_website_configuration" {
  bucket = aws_s3_bucket.clams_frontend_static_content_bucket.bucket
  index_document {
    suffix = "index.html"
  }
  error_document {
    key = "index.html"
  }
}

resource "aws_s3_bucket_public_access_block" "clams_frontend_bucket_public_access_block" {
  bucket = aws_s3_bucket.clams_frontend_static_content_bucket.id

  block_public_acls        = true
  block_public_policy      = true
  ignore_public_acls       = true
  restrict_public_buckets  = true
}

data "aws_iam_policy_document" "clams_frontend_bucket_policy" {

  statement {
    sid = "PublicReadForGetBucketObjects"
    principals {
      type        = "AWS"
      identifiers = [var.origin_access_identity]
    }
    actions = [
      "s3:GetObject",
    ]
    resources = [
      "arn:aws:s3:::${var.frontend_domain}/*",
    ]
  }
}