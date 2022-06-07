module "dynamo" {
  source         = "./modules/dynamo"
  environment    = var.environment
  region         = var.region
  account_number = var.account_number
  contact        = var.contact
  product        = var.product
  orchestration  = var.orchestration
  attendees_table_name = var.attendees_table_name
}

module "sqs" {
  source           = "./modules/sqs"
  environment      = var.environment
  region           = var.region
  account_number   = var.account_number
  contact          = var.contact
  product          = var.product
  orchestration    = var.orchestration
  input_queue_name = var.input_queue_name
}

module "lambda" {
  source               = "./modules/lambda"
  environment          = var.environment
  region               = var.region
  account_number       = var.account_number
  contact              = var.contact
  product              = var.product
  orchestration        = var.orchestration
  distribution_bucket  = var.distribution_bucket
  attendees_table_arn  = module.dynamo.attendees_table_arn
  attendees_table_name = var.attendees_table_name
  input_queue_arn      = module.sqs.attendee_input_queue_arn
  input_queue_name     = var.input_queue_name
}

module "s3" {
  source = "./modules/s3"

  environment = var.environment
  contact = var.contact
  product = var.product
  orchestration = var.orchestration
  origin_access_identity = module.cloudfront.origin_access_identity.iam_arn
  frontend_domain = var.frontend_domain
}

module "route53" {
  source = "./modules/route53"

  cloudfront_domain_name = module.cloudfront.cloudfront_domain_name
  certificate_domain = var.certificate_domain
  frontend_domain  = var.frontend_domain
}

module "cloudfront" {
  source = "./modules/cloudfront"

  environment = var.environment
  contact = var.contact
  product = var.product
  orchestration = var.orchestration

  acm_certificate_arn = data.aws_acm_certificate.clams_cert.arn
  origin_domain_name = module.s3.bucket_regional_domain_name
  route53_zone_id = module.route53.aws_route53_zone_id
  certificate_domain = var.certificate_domain
  frontend_domain  = var.frontend_domain
}
