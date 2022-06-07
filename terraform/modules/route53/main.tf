resource "aws_route53_record" "clams_frontend_a_record" {
  zone_id = data.aws_route53_zone.clams_r53_zone.id
  name    = "clams"
  type    = "A"

  alias {
    name                   = var.cloudfront_domain_name
    zone_id                = var.cloudfront_hosted_zone_id
    evaluate_target_health = false
  }
}

data "aws_route53_zone" "clams_r53_zone" {
  name         = var.certificate_domain
  private_zone = false
}
