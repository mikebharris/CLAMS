resource "aws_sqs_queue" "attendee_input_queue" {
  name = "${var.product}-attendee-input-queue"

  redrive_policy = jsonencode({
    "deadLetterTargetArn" = aws_sqs_queue.attendee_input_dlq.arn,
    "maxReceiveCount"     = var.receive_count
  })

  tags = {
    Name          = "${var.product}.${var.environment}.sqs.attendee_input_queue"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    SubProduct    = var.sub_product
    CostCode      = var.cost_code
    Orchestration = var.orchestration
    Description   = "SQS Queue for holding attendees"
  }
}

resource "aws_sqs_queue" "attendee_input_dlq" {
  name                      = "${var.product}-attendee-input-queue-DLQ"
  message_retention_seconds = var.dlq_retention_period

  tags = {
    Name          = "${var.product}.${var.environment}.sqs.attendee_input_queue_DLQ"
    Contact       = var.contact
    Environment   = var.environment
    Product       = var.product
    SubProduct    = var.sub_product
    CostCode      = var.cost_code
    Orchestration = var.orchestration
    Description   = "SQS dead letter queue for attendees that could not be processed by the eHAMS Attendee Librarian"
  }
}

resource "aws_sqs_queue_policy" "attendee_input_dlq_policy" {
  queue_url = aws_sqs_queue.attendee_input_dlq.id
  policy    = jsonencode({
    Version   = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowAttendeeInputQueueToSend"
        Effect    = "Allow"
        Principal = "*"
        Action    = "sqs:SendMessage"
        Resource  = aws_sqs_queue.attendee_input_dlq.arn
        Condition = {
          ArnEquals = { "AWS:SourceArn" = aws_sqs_queue.attendee_input_queue.arn }
        }
      }
    ]
  })
}
