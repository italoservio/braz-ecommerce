resource "aws_sqs_queue" "user_created_queue" {
  name                       = "${local.resource_prefix}-user-created-queue"
  max_message_size           = 2048
  visibility_timeout_seconds = 60
}
