variable environment {}
variable region {}
variable account_number {}
variable contact {}
variable product {}
variable sub_product {}
variable cost_code {}
variable orchestration {}

variable distribution_bucket {}
variable cusoon_results_table_arn {}
variable cusoon_email_sends_arn {}
variable cusoon_supplier_results_queue_arn {}
variable ses_domain{}
variable send_emails{}
variable to_email_override{}
variable frontend_url{}
variable enabled_error_codes{}
variable cusoon_results_datalake_bucket_arn {}

variable "cors_configuration" {
  type = any
  default = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }
}