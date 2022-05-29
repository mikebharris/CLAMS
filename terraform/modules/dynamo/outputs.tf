output "attendees_table_arn" {
  value = aws_dynamodb_table.attendees_table.arn
}

output "attendees_table_name" {
  value = aws_dynamodb_table.attendees_table.name
}