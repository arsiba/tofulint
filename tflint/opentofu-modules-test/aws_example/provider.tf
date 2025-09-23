# run `tofulint --init` before testing to install the aws module

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }

  required_version = "~> 1.1.9"
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-1"
}