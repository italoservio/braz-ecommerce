variable "aws_access_key_id" {
  type        = string
  description = "AWS Access Key Id"
}

variable "aws_secret_access_key" {
  type        = string
  description = "AWS Secret Access Key"
}

variable "aws_region" {
  type        = string
  description = "AWS region"
}

variable "aws_account_id" {
  type        = string
  description = "AWS Account Id"
}

variable "aws_endpoint" {
  type        = string
  description = "AWS custom endpoint"
}

variable "environment" {
  type        = string
  description = "Current environment"
}
