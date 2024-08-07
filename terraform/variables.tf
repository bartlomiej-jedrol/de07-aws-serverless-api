// Terraform variables.
variable "region" {
  type    = string
  default = "eu-central-1"
}

variable "lambda_iam_role" {
  type    = string
  default = "de07-lambda-role"
}

variable "lambda_inline_policy" {
  type    = string
  default = "de07-lambda-inline-policy"
}

variable "lambda_function_name" {
  type    = string
  default = "de07-lambda"
}

variable "api_gateway_name" {
  type    = string
  default = "de07-api-gateway"
}

variable "dynamodb_table_name" {
  type    = string
  default = "de07-user"
}
