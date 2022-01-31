provider "aws" {
  region = "ap-northeast-1"
}

terraform {
  required_providers {
    aws = {
      version = "~> 3.38.0"
    }
  }
}