terraform {
  backend "local" {
    path = "../../.state/terraform.tfstate"
  }
}

provider "aws" {
  access_key = var.aws_access_key_id
  secret_key = var.aws_secret_access_key

  region                      = var.aws_region
  s3_use_path_style           = true
  skip_requesting_account_id  = true
  skip_metadata_api_check     = true
  skip_credentials_validation = true

  dynamic "endpoints" {
    for_each = var.environment == "local" ? [1] : []

    content {
      apigateway     = var.aws_endpoint
      cloudformation = var.aws_endpoint
      cloudwatch     = var.aws_endpoint
      dynamodb       = var.aws_endpoint
      es             = var.aws_endpoint
      firehose       = var.aws_endpoint
      iam            = var.aws_endpoint
      kinesis        = var.aws_endpoint
      lambda         = var.aws_endpoint
      route53        = var.aws_endpoint
      redshift       = var.aws_endpoint
      s3             = var.aws_endpoint
      secretsmanager = var.aws_endpoint
      ses            = var.aws_endpoint
      sns            = var.aws_endpoint
      sqs            = var.aws_endpoint
      ssm            = var.aws_endpoint
      stepfunctions  = var.aws_endpoint
      sts            = var.aws_endpoint
    }
  }
}
