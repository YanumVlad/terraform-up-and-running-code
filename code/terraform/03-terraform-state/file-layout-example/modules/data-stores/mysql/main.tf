terraform {
  # Require any 0.12.x version of Terraform
  required_version = ">= 0.12, <= 0.13"

  # Partial configuration. The rest will be filled in by Terragrunt.
  backend "s3" {}
}
provider "aws" {
  region = "us-east-2"

  # Allow any 2.x version of the AWS provider
  version = "~> 2.0"
}

resource "aws_db_instance" "example" {

  identifier        = var.instance_name
  engine            = "mysql"
  allocated_storage = 10
  instance_class    = "db.t2.micro"

  username = "admin"
  password = var.db_password

  name                = var.db_name
  skip_final_snapshot = true

}

