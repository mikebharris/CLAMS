output "attendee_input_queue_arn" {
  value = aws_sqs_queue.attendee_input_queue.arn
}

output "attendee_input_queue_name" {
  value = aws_sqs_queue.attendee_input_queue.name
}