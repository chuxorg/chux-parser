terraform {
  backend "s3" {
    bucket = "chux-terraform-state"
    key    = "chux-lambda-terraform.tfstate"
    region = "us-east-1"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "us-east-1" # Change this to your desired AWS region
}

locals {
  function_name = "chux-ecs-parse"
}

resource "aws_ecs_cluster" "chux_cluster" {
  name = "chux-cluster"
}

resource "aws_ecs_task_definition" "chux_task_definition" {
  family                = "chux-task-family"
  network_mode          = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                   = "256"
  memory                = "512"
  execution_role_arn    = aws_iam_role.ecs_execution_role.arn
  container_definitions = jsonencode([{
    name  = "chux-container"
    image = "${aws_ecr_repository.chux_lambda_parser.repository_url}:latest"
    essential = true
    environment = [
      { name = "AWS_BUCKET", value = "chux-crawler" },
      { name = "AWS_SOURCE_BUCKET", value = "chux-crawler" },
      { name = "LOG_FILE_NAME", value = "chux-cprs" },
      { name = "MONGO_URI", value = "mongodb+srv://%s:%s@chux-mongo-cluster.4mvs7.mongodb.net/" },
      { name = "MONGO_USER_NAME", value = "username" },
      { name = "MONGO_PASSWORD", value = "password" },
      { name = "MONGO_DATABASE", value = "chux-cprs" },
      { name = "ENVIRONMENT", value = "development" },
      { name = "TARGET", value = "Fargate" },
    ]
  }])
}

resource "aws_ecs_service" "chux_service" {
  name            = local.service_name
  cluster         = aws_ecs_cluster.chux_cluster.id
  task_definition = aws_ecs_task_definition.chux_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = ["subnet-009f7d01c00791a01"]
    security_groups  = [aws_security_group.lambda_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.chux_tg.arn
    container_name   = local.container_name
    container_port   = 80
  }

  depends_on = [aws_lb_listener.chux_listener]
}


resource "aws_security_group" "lambda_sg" {
  name        = "${local.function_name}_sg"
  description = "Security group for Lambda function to access the internet"
  vpc_id      = "vpc-0d29c91c33cb0acd7"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"] # VPC's CIDR block
  }
}

resource "aws_iam_role" "ecs_execution_role" {
  name = "${local.function_name}_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  role       = aws_iam_role.ecs_execution_role.id
}

resource "aws_iam_role_policy_attachment" "s3_secretsmanager_cloudwatch" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
  role       = aws_iam_role.ecs_execution_role.id
}

resource "aws_iam_role_policy_attachment" "secretsmanager_policy" {
  policy_arn = "arn:aws:iam::aws:policy/SecretsManagerReadWrite"
  role       = aws_iam_role.ecs_execution_role.id
}

resource "aws_iam_role_policy_attachment" "cloudwatch_policy" {
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchFullAccess"
  role       = aws_iam_role.ecs_execution_role.id
}

resource "aws_ecr_repository" "chux_lambda_parser" {
  name = "chux-lambda-parser"
}

