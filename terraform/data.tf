data "aws_acm_certificate" "clams_cert" {
  domain   = var.certificate_domain
  statuses = ["ISSUED"]
}