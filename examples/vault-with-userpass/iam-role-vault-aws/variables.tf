variable "name" {
  description = "The name of the role to create."
  type        = string
}

variable "description" {
  description = "The description of the role."
  type        = string
}

variable "path" {
  description = "The AWS path under which to create the role."
  default     = "/"
  type        = string
}

variable "assume_policy_principal_arn" {
  description = "The ARN of the user used to assume this role and generate credentials. Should be the ARN of the root user configured into the Vault AWS engine."
  type = string
}

variable "policies" {
  description = "A map of policies to attach to the role, indexed by their name."
  type        = map(string)
}

variable "default_sts_ttl" {
  description = "The default Vault TTL seconds for the STS credentials."
  type        = number
  default     = 7200
}

variable "max_sts_ttl" {
  description = "The max Vault TTL seconds for the STS credentials."
  type        = number
  default     = 7200
}

variable "vault_mount_path" {
  description = "The path of the aws engine in Vault to create the role."
  default     = "aws"
  type        = string
}


