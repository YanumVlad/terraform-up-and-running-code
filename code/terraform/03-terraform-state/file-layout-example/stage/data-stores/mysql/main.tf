terraform {

  required_version = ">= 0.12, <= 0.13"

  backend "s3" {
    # Replace this with your bucket name!
    bucket = "ulny-s3"
    key    = "stage/data-stores/mysql/terraform.tfstate"
    region = "us-east-2"

    # Replace this with your DynamoDB table name!
    dynamodb_table = "dynamo-lock-table"
    encrypt        = true
  }
}
provider "aws" {
  region = "us-east-2"

  # Allow any 2.x version of the AWS provider
  version = "~> 2.0"
}

resource "aws_db_instance" "example" {
  identifier_prefix = "terraform-up-and-running"
  engine            = "mysql"
  allocated_storage = 10
  instance_class    = "db.t2.micro"

  username = "admin"

  name                = var.db_name
  skip_final_snapshot = true

  password = var.db_password
}

