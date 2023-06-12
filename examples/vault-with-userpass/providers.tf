terraform {
  required_providers {
    vault = "~> 3"
    aws   = ">= 5, < 6"
  }
}

provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile
  default_tags {
    tags = {
      "ManagedBy" = "Terraform"
      "Project"   = "https://github.com/thomasgouveia/vault-aws-credentials-helper"
    }
  }
}

provider "vault" {
  address = "http://localhost:8200"
  # Do not ever do this in production.
  # This is just for development purposes.
  token = "root"
}
