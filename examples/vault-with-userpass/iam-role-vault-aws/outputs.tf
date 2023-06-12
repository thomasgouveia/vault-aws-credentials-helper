output "role_arn" {
  value = aws_iam_role.this.arn
}

output "role_name" {
  value = var.name
}
