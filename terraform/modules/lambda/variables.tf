variable environment {}
variable region {}
variable account_number {}
variable contact {}
variable product {}
variable orchestration {}
variable distribution_bucket {}
variable "attendees_table_arn" {}
variable "attendees_table_name" {}
variable "attendees_queue_arn" {}
variable "attendees_queue_name" {}

variable "signups_queue_arn" {}

variable "cors_configuration" {
  type = any
  default = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }
}

variable "db_host" {}
variable "db_name" {}
variable "db_password" {}
variable "db_username" {}