terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 4.46"
    }
  }
}

provider "aws" {
  region = var.region
}

terraform {
  backend "s3" {}
}
