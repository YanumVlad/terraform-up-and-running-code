terraform {

  required_version = ">= 0.12, <= 0.13"

  backend "s3" {
    bucket = "ulny-s3"
    key    = "stage/data-stores/mysql/terraform.tfstate"
    region = "us-east-2"

    dynamodb_table = "dynamo-lock-table"
    encrypt        = true
  }
}
provider "aws" {
  region = "us-east-2"

  # Allow any 2.x version of the AWS provider
  version = "~> 2.0"
}

module "mysql" {

  source = "../../../modules/data-stores/mysql"

  instance_name = var.cluster_name
  db_name       = var.db_name
  db_password   = var.db_password
}
