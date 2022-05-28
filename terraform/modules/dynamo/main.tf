resource "aws_dynamodb_table" "attendees_table" {
  name         = "mbh-attendees-datastore"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Code"

  attribute {
    name = "Code"
    type = "S"
  }

  tags = {
    Name          = "${var.product}.${var.environment}.attendees_datastore"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    SubProduct    = var.sub_product
    CostCode      = var.cost_code
    Orchestration = var.orchestration
    Description   = "Store of attendees from BAMS"
  }
  point_in_time_recovery {
    enabled = false
  }
}
