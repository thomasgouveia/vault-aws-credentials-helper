# Create the policy allowing the specified user (through its arn) to assume this role.
# It will be used by Vault to automatically, through the root IAM user, assume the role 
# and retrieve credentials from STS.
data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "AWS"
      identifiers = [var.assume_policy_principal_arn]
    }
  }
}

# Create the role in the AWS account
resource "aws_iam_role" "this" {
  name               = var.name
  description        = var.description
  path               = var.path
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json

  # Dynamically generate the inline policies for this role.
  dynamic "inline_policy" {
    for_each = var.policies
    content {
      name   = inline_policy.key
      policy = inline_policy.value
    }
  }
}

# Create the associated role in the Vault, so that
# users will be able to perform credentials generation.
resource "vault_aws_secret_backend_role" "this" {
  name            = var.name
  default_sts_ttl = var.default_sts_ttl
  max_sts_ttl     = var.max_sts_ttl
  backend         = var.vault_mount_path
  role_arns       = [aws_iam_role.this.arn]
  credential_type = "assumed_role"
}
