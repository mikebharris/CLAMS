terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 4.46"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

terraform {
  backend "s3" {}
}
