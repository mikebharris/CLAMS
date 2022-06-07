output "aws_route53_zone_id" {
  value = data.aws_route53_zone.clams_r53_zone.id
}

output "route53_hosted_zone_domain" {
  value = data.aws_route53_zone.clams_r53_zone
}