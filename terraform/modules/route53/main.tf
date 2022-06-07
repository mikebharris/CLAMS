resource "aws_route53_record" "clams_frontend_a_record" {
  zone_id = data.aws_route53_zone.clams_r53_zone.id
  name    = var.frontend_domain
  type    = "A"

  alias {
    name                   = var.cloudfront_domain_name
    //zone_id value is static
    zone_id                = "Z0640145NHPXT7F33V6B"
    evaluate_target_health = false
  }
}

data "aws_route53_zone" "clams_r53_zone" {
  name         = var.certificate_domain
  private_zone = false
}
