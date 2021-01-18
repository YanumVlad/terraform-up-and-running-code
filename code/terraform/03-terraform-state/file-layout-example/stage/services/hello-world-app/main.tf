provider "aws" {
  region = "us-east-2"

  # Allow any 2.x version of the AWS provider
  version = "~> 2.0"
}

module "hello_world_app" {
  source = "github.com/YanumVlad/terraform-modules.git//modules/services/hello-world-app?ref=v0.0.5"

  server_text = "New server text"
  environment = "stage"
  bucket      = "ulny-s3"
  key         = "stage/services/webserver-cluster/terraform.tfstate"

  instance_type      = "t2.micro"
  min_size           = 2
  max_size           = 2
  enable_autoscaling = false
}
