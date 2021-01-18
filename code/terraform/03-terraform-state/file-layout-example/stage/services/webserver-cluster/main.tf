terraform {
  required_version = ">= 0.12, <= 0.13"

  backend "s3" {
    # Replace this with your bucket name!
    bucket = "ulny-s3"
    key    = "stage/services/webserver-cluster/terraform.tfstate"
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

  cluster_name                    = var.cluster_name
  db_remote_state_bucket          = var.db_remote_state_bucket
  db_remote_state_key             = var.db_remote_state_key
  instance_type                   = var.instance_type
  min_size                        = var.min_size
  max_size                        = var.max_size
  enable_autoscaling              = false
  give_neo_cloudwatch_full_access = true
  ami                             = "ami-0c55b159cbfafe1f0"
  server_text                     = "New server text IIHF"


  custom_tags = {
    Owner       = "team-foo"
    DeployedBy  = "terraform"
    Environment = "staging"
  }
}

variable "names" {
  description = "A list of names"
  type        = list(string)
  default     = ["neo", "trinity", "morpheus"]
}

output "upper_names" {
  value = [for name in var.names : upper(name) if length(name) < 5]

}

variable "hero_thousand_faces" {
  description = "map"
  type        = map(string)
  default = {
    neo      = "hero"
    trinity  = "love interest"
    morpheus = "mentor"
  }
}

output "bios" {
  value = [for name, role in var.hero_thousand_faces : "${name} is the ${role}"]
}

output "upper_roles" {
  value = { for name, role in var.hero_thousand_faces : upper(name) => upper(role) }
}

output "for_directive" {
  value = <<EOF
%{~for name in var.names}
  ${name}
%{~endfor}
EOF
}

# variable "name" {
#   description = "A name to render"
#   type        = string
# }

# output "if_else_directive" {
#   value = "Hello, %{if var.name != ""}${var.name}%{else}(unnamed)%{endif}"
# }

#resource "aws_iam_user" "existing_user" {
#  # You should change this to the username of an IAM user that already
#  # exists so you can practice using the terraform import command
#  name = "developer1"
#}
