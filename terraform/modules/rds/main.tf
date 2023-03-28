resource "aws_db_instance" "hacktionlab_signups_database" {
  allocated_storage    = 10
  identifier           = "hacktionlab"
  db_name              = "hacktionlab"
  engine               = "postgres"
  availability_zone    = "us-east-1a"
  instance_class       = "db.t4g.micro"
  username             = data.aws_ssm_parameter.db_username.value
  password             = data.aws_ssm_parameter.db_password.value
  skip_final_snapshot  = true
  publicly_accessible  = true
  parameter_group_name = aws_db_parameter_group.default.name
}

resource "aws_db_parameter_group" "default" {
  name   = "postgres-connection-logging"
  family = "postgres14"

  parameter {
    name  = "log_connections"
    value = "1"
  }
}

