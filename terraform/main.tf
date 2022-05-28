module "dynamo" {
  source         = "./modules/dynamo"
  environment    = var.environment
  region         = var.region
  account_number = var.account_number
  contact        = var.contact
  product        = var.product
  sub_product    = var.sub_product
  cost_code      = var.cost_code
  orchestration  = var.orchestration
}

module "sqs" {
  source         = "./modules/sqs"
  environment    = var.environment
  region         = var.region
  account_number = var.account_number
  contact        = var.contact
  product        = var.product
  sub_product    = var.sub_product
  cost_code      = var.cost_code
  orchestration  = var.orchestration
}
#
#module "lambda" {
#  source = "./modules/lambda"
#  environment    = var.environment
#  region         = var.region
#  account_number = var.account_number
#  contact        = var.contact
#  product        = var.product
#  sub_product    = var.sub_product
#  cost_code      = var.cost_code
#  orchestration  = var.orchestration
#  distribution_bucket = var.distribution_bucket
#  cusoon_results_table_arn = module.dynamo.cusoon_results_table_arn
#  cusoon_email_sends_arn = module.dynamo.cusoon_email_sends_arn
#  cusoon_supplier_results_queue_arn = module.sqs.cusoon_supplier_results_queue_arn
#  ses_domain = var.ses_domain
#  send_emails = var.send_emails
#  to_email_override = var.to_email_override
#  frontend_url = var.frontend_url
#  enabled_error_codes = var.enabled_error_codes
#  cusoon_results_datalake_bucket_arn = module.s3.cusoon_results_datalake_arn
#}
#
#module "s3" {
#  source = "./modules/s3"
#  environment = var.environment
#  region = var.region
#  account_number = var.account_number
#  contact = var.contact
#  product = var.product
#  sub_product = var.sub_product
#  cost_code = var.cost_code
#  orchestration = var.orchestration
#}