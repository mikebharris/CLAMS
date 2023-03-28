output "rds_database_password" {
  value = data.aws_ssm_parameter.db_password.value
}

output "rds_database_username" {
  value = data.aws_ssm_parameter.db_username.value
}

output "rds_database_host" {
  value = aws_db_instance.hacktionlab_signups_database.address
}

output "rds_database_name" {
  value = aws_db_instance.hacktionlab_signups_database.db_name
}