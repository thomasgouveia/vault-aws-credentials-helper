# Basic example

This folder contains a basic example that can be used to test the basic functionality of the `vault-aws-credentials-helper`. 

## Prerequisites

- An AWS account configured on your laptop
- [Terraform](https://developer.hashicorp.com/terraform/downloads?product_intent=terraform) installed on your laptop
- [Vault](https://developer.hashicorp.com/vault/downloads?product_intent=vault) installed on your laptop

In a first terminal, start vault in dev mode:

```bash
vault server -dev -dev-root-token-id root
```

Create a file `terraform.tfvars` on this folder on your laptop with the following content : 

```hcl
aws_profile = "<YOUR_AWS_PROFILE>"
aws_region  = "<YOUR_AWS_REGION>"
```

Now, you can use Terraform in another terminal to deploy and configure all the resources needed:

```bash
terraform init
terraform apply
```

If everything is ok, you should have the following output:

```bash
Apply complete! Resources: 17 added, 0 changed, 0 destroyed.
```

Now, we can test that everything is working fine. The Terraform configuration generates two users in Vault:

- `john.smith / mysuperpassword` : part of a group `devteam`, **allowed to generate `developer` credentials**
- `sam.wellington / myanotherpassword`: part of a group `qateam`, **not allowed to generate `developer` credentials**

Also, Terraform configures a `developer` IAM role on AWS allowing full access to EC2 in the `eu-west-3` region **ONLY**.

We can test that we can generate a set of credentials with the user `john.smith`:

```bash
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=$(vault login -format=json -method=userpass username=john.smith password=mysuperpassword | jq -r '.auth.client_token')

# Try to get AWS credentials, it should work as this user is part of the devteam
vault write -format=json aws/sts/developer ttl=15m
```

Now, ensure that the `sam.wellington` user can't generate credentials:

```bash
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=$(vault login -format=json -method=userpass username=sam.enerv password=myanotherpassword | jq -r '.auth.client_token')

# Try to get AWS credentials, we should have a 403 permission denied as this user is part of the qateam and this team
# can't generate AWS credentials
vault write -format=json aws/sts/developer ttl=15m
```

## Use with `vault-aws-credentials-helper`

Now that your setup is ready, you can simply use the `vault-aws-credentials-helper` to issue credentials :

```bash
vault-aws-credentials-helper --vault.addr http://localhost:8200 \
    --vault.auth-method userpass \
    --userpass.username john.smith \
    --userpass.password mysuperpassword \
    --aws.role developer \
    --aws.ttl 30m
```

You should have the following output:

```json
{
 "Version": 1,
 "AccessKeyId": "<REDACTED>",
 "SecretAccessKey": "<REDACTED>",
 "SessionToken": "<REDACTED>",
 "Expiration": "2023-06-11T22:43:22.017Z"
}
```

Please refer to the [project README](../../README.md) to configure your AWS CLI to use the `vault-aws-credentials-helper` to generate credentials.

## Clean up

To clean up every resources created by this configuration, simply run:

```bash
terraform destroy
```
