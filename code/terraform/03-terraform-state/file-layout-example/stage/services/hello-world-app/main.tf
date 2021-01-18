provider "aws" {
  region = "us-east-2"

  # Allow any 2.x version of the AWS provider
  version = "~> 2.0"
}

module "hello_world_app" {
  source = "github.com/YanumVlad/terraform-modules.git//modules/services/hello-world-app?ref=v0.0.7"

  server_text            = "New server text"
  environment            = "stage"
  db_remote_state_bucket = "ulny-s3"
  db_remote_state_key    = "stage/data-stores/mysql/terraform.tfstate"
  enable_autoscaling     = true

  instance_type = "t2.micro"
  min_size      = 2
  max_size      = 2
}
