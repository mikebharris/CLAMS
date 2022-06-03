resource "aws_dynamodb_table" "attendees_table" {
  name         = var.attendees_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "AuthCode"

  attribute {
    name = "AuthCode"
    type = "S"
  }

  tags = {
    Name          = "${var.product}.${var.environment}.attendees_datastore"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    Orchestration = var.orchestration
    Description   = "Store of attendees from BAMS"
  }
  point_in_time_recovery {
    enabled = false
  }
}
