terraform {

  required_version = ">= 0.12, <= 0.13"
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

