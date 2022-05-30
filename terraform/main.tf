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
  attendees_table_name = var.attendees_table_name
}

module "sqs" {
  source           = "./modules/sqs"
  environment      = var.environment
  region           = var.region
  account_number   = var.account_number
  contact          = var.contact
  product          = var.product
  sub_product      = var.sub_product
  cost_code        = var.cost_code
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
  sub_product          = var.sub_product
  cost_code            = var.cost_code
  orchestration        = var.orchestration
  distribution_bucket  = var.distribution_bucket
  attendees_table_arn  = module.dynamo.attendees_table_arn
  attendees_table_name = var.attendees_table_name
  input_queue_arn      = module.sqs.attendee_input_queue_arn
  input_queue_name     = var.input_queue_name
}
