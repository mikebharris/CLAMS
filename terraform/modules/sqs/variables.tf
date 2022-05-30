variable environment {}
variable region {}
variable account_number {}
variable contact {}
variable product {}
variable sub_product {}
variable cost_code {}
variable orchestration {}
variable input_queue_name {}

variable "receive_count" {
  description = "The number of times that a message can be retrieved before being moved to the dead-letter queue"
  type = number
  default = 3
}

variable "dlq_retention_period" {
  description = "Time (in seconds) that messages will remain in queue before being purged"
  type = number
  default = 1209600 #14 days, default is 4 days max is 14 days.
}