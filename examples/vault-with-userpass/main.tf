data "aws_caller_identity" "current" {}

# Create the IAM user that will be used by Vault to manage dynamic IAM users.
# This requires that the user applying this Terraform configuration to have sufficient
# permissions to create the user and attach a policy to it.
data "aws_iam_policy_document" "vault_iam_manager" {
  statement {
    effect = "Allow"
    # Important to restrict only to users managed by Vault
    resources = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/vault-*"]
    actions = [
      "iam:AttachUserPolicy",
      "iam:CreateAccessKey",
      "iam:CreateUser",
      "iam:DeleteAccessKey",
      "iam:DeleteUser",
      "iam:DeleteUserPolicy",
      "iam:DetachUserPolicy",
      "iam:GetUser",
      "iam:ListAccessKeys",
      "iam:ListAttachedUserPolicies",
      "iam:ListGroupsForUser",
      "iam:ListUserPolicies",
      "iam:PutUserPolicy",
      "iam:AddUserToGroup",
      "iam:RemoveUserFromGroup"
    ]
  }
}

resource "aws_iam_user" "vault_iam_manager" {
  name = "vault-iam-manager"
}

resource "aws_iam_access_key" "vault_iam_manager" {
  user = aws_iam_user.vault_iam_manager.name
}

resource "aws_iam_policy" "vault_iam_manager" {
  name        = "VaultIAMManagerPolicy"
  description = "The policy that should be used by the Vault IAM Manager user."
  policy      = data.aws_iam_policy_document.vault_iam_manager.json
}

resource "aws_iam_user_policy_attachment" "vault_iam_manager" {
  user       = aws_iam_user.vault_iam_manager.name
  policy_arn = aws_iam_policy.vault_iam_manager.arn
}

# Configure the Vault AWS backend engine with the freshly created user
# Do not forget to configure the max_least_ttl_seconds in order to ensure
# that generated credentials will be short lived.
resource "vault_aws_secret_backend" "aws" {
  description           = "Dynamic credentials for AWS"
  region                = "eu-west-3"
  access_key            = aws_iam_access_key.vault_iam_manager.id
  secret_key            = aws_iam_access_key.vault_iam_manager.secret
  max_lease_ttl_seconds = 7200 # 2hours
}

# Create a role on AWS for developers. For the example here, 
# we assume that the developers will have full access to ec2 in
# the eu-west-3 region only.
data "aws_iam_policy_document" "developer" {
  statement {
    effect    = "Allow"
    actions   = ["ec2:*"]
    resources = ["*"]

    condition {
      test     = "StringEquals"
      variable = "aws:RequestedRegion"
      values   = ["eu-west-3"]
    }
  }
}

# This module abstracts the creation of the role on AWS and its mapping in Vault.
module "role_developer" {
  source      = "./iam-role-vault-aws"
  name        = "developer"
  description = "The role each developer should assume to access cloud resources."

  # ARN of the user configured into the Vault AWS engine
  assume_policy_principal_arn = aws_iam_user.vault_iam_manager.arn
  # IAM Policies to attach to the role
  policies = {
    DeveloperPolicy = data.aws_iam_policy_document.developer.json
  }
}

# Create a policy in Vault that will be attached to our
# developer group so every users in the group will be able 
# to generate dynamic AWS credentials for the developer role on AWS
resource "vault_policy" "developer" {
  name   = "developer"
  policy = <<EOT
    path "${vault_aws_secret_backend.aws.path}/sts/${module.role_developer.role_name}" {
      capabilities = ["create", "read", "update"]
    }
  EOT
}

# Configure a stub user in Vault through identity.
# We assume that the user is a developer, it will be added to the developer group.
resource "vault_identity_entity" "john" {
  name = "John Smith"
  metadata = {
    email = "john.smith@gmail.com"
  }
}

# Configure another stub user in Vault through identity.
# We assume that the user is a QA member, it will be added to the QA group.
resource "vault_identity_entity" "sam" {
  name = "Sam Wellington"
  metadata = {
    email = "sam.wellington@gmail.com"
  }
}

# Create an internal group to manage all our fake developer team
resource "vault_identity_group" "devteam" {
  name              = "devteam"
  type              = "internal"
  member_entity_ids = [vault_identity_entity.john.id]
  policies = [
    vault_policy.developer.name
  ]
}

# Create an internal group to manage all our fake qa team
resource "vault_identity_group" "qateam" {
  name              = "qateam"
  type              = "internal"
  member_entity_ids = [vault_identity_entity.sam.id]
  # do not attach any policies, to illustrate the management of roles / policies
  policies = []
}

# Enable an authentication method, for instance, userpass
resource "vault_auth_backend" "userpass" {
  type = "userpass"
  path = "userpass"
}

# Create the John Smith user into the userpass auth backend
resource "vault_generic_endpoint" "john_smith_user" {
  path                 = "auth/${vault_auth_backend.userpass.path}/users/john.smith"
  ignore_absent_fields = true
  data_json = jsonencode({
    password = "mysuperpassword"
  })
}

# Create the Sam Wellington user into the userpass auth backend
resource "vault_generic_endpoint" "sam_wellington_user" {
  path                 = "auth/${vault_auth_backend.userpass.path}/users/sam.wellington"
  ignore_absent_fields = true
  data_json = jsonencode({
    password = "myanotherpassword"
  })
}

# Map the userpass auth for John Smith to its identity through Vault
resource "vault_identity_entity_alias" "john_smith" {
  name           = "john.smith"
  mount_accessor = vault_auth_backend.userpass.accessor
  canonical_id   = vault_identity_entity.john.id
}

# Map the userpass auth for Sam Wellington to its identity through Vault
resource "vault_identity_entity_alias" "sam_wellington" {
  name           = "sam.wellington"
  mount_accessor = vault_auth_backend.userpass.accessor
  canonical_id   = vault_identity_entity.sam.id
}
