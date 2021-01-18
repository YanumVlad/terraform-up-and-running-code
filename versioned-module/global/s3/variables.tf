variable "bucket_name" {
  description = "The name of the S3 bucket. Must be globally unique."
  default     = "ulny-s3"
  type        = string
}

variable "table_name" {
  description = "The name of the DynamoDB table. Must be unique in this AWS account."
  default     = "dynamo-lock-table"
  type        = string
}
