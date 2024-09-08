terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket  = "bj-terraform-states"
    key     = "state-de07-aws-serverless-api/terraform.tfstate"
    region  = "eu-central-1"
    encrypt = true
  }
}

provider "aws" {
  region = "eu-central-1"
}

# Lambda
resource "aws_iam_role" "lambda_iam_role" {
  name = var.lambda_iam_role
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })

  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
}

resource "aws_lambda_function" "lambda_function" {
  function_name = var.lambda_function_name
  handler       = "main"
  runtime       = "provided.al2023"
  timeout       = 10
  filename      = "../build/main.zip"

  role = aws_iam_role.lambda_iam_role.arn
}

resource "aws_iam_role_policy" "lambda_inline_policy" {
  name = var.lambda_inline_policy
  role = aws_iam_role.lambda_iam_role.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Scan",
        ]
        Resource = [aws_dynamodb_table.dynamodb_table.arn]
      },
    ],
  })
}

# API Gateway
resource "aws_api_gateway_rest_api" "api_gateway" {
  name = var.api_gateway_name
}

# DynamoDB
resource "aws_dynamodb_table" "dynamodb_table" {
  name         = var.dynamodb_table_name
  hash_key     = "email"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "email"
    type = "S"
  }

  attribute {
    name = "firstName"
    type = "S"
  }

  attribute {
    name = "lastName"
    type = "S"
  }

  attribute {
    name = "age"
    type = "N"
  }

  global_secondary_index {
    name            = "FirstNameIndex"
    hash_key        = "firstName"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "LastNameIndex"
    hash_key        = "lastName"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "AgeIndex"
    hash_key        = "age"
    projection_type = "ALL"
  }
}

# CloudWatch
resource "aws_cloudwatch_log_group" "cloud_watch_group" {
  name = "/aws/apigateway/${var.api_gateway_name}"
}
