terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"

    }
  }

  required_version = "~> 1.1.9"
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-1"
}