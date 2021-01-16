terraform {
  required_version = ">= 0.12, <= 0.13"

  backend "s3" {
    # Replace this with your bucket name!
    bucket = "ulny-s3"
    key    = "prod/services/webserver-cluster/terraform.tfstate"
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
module "webserver_cluster" {
  source = "../../../modules/services/webserver-cluster"

  cluster_name           = "webservers-prod"
  db_remote_state_bucket = "ulny-s3"
  db_remote_state_key    = "prod/data-stores/mysql/terraform.tfstate"
  instance_type          = "t2.micro"
  min_size               = 3
  max_size               = 10
}
