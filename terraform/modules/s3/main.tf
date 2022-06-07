resource "aws_s3_bucket" "clams_frontend_static_content_bucket" {
  bucket        = var.frontend_domain
  acl           = "public-read"
  force_destroy = true

  website {
    index_document = "index.html"
    error_document = "index.html"
  }

  policy = data.aws_iam_policy_document.clams_frontend_bucket_policy.json

  tags = {
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    Orchestration = var.orchestration
    Description   = "S3 bucket for CLAMS frontend static content in ${var.environment} environment."
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