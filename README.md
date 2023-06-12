# Vault AWS Credentials Helper

`vault-aws-credentials-helper` is a simple CLI tool that provides the ability to generate dynamic AWS credentials using the [HashiCorp Vault's AWS secret engine](https://developer.hashicorp.com/vault/docs/secrets/aws) and integrate seamlessly with the AWS CLI to retrieve dynamic and short-lived credentials. **Short-lived credentials enforce security and reduce the risk of a credential leak or corruption**. 

To learn more about how the credentials are provided to the AWS CLI, please refer to the official AWS documentation: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html.

## Supported AWS credentials

**Currently, the tool supports only the generation of [STS AssumeRole](https://developer.hashicorp.com/vault/docs/secrets/aws#sts-assumerole) credentials.**

The support for [STS Federation Tokens](https://developer.hashicorp.com/vault/docs/secrets/aws#sts-federation-tokens) and [IAM Users](https://developer.hashicorp.com/vault/docs/secrets/aws#iam_user) will be available in a future version.


## Supported Vault authentication methods

To be able to use a secret engine in Vault, you must be able to retrieve an authentication token from it. [Vault supports a lot of authentication methods](https://developer.hashicorp.com/vault/docs/auth), and below the list of the supported authentication method in the tool currently:

- [Token](https://developer.hashicorp.com/vault/docs/auth/token)
- [Userpass](https://developer.hashicorp.com/vault/docs/auth/userpass)

## Install

You can install `vault-aws-credentials-helper` easily using `go install`: 

```bash
go install github.com/thomasgouveia/vault-aws-credentials-helper@v0.1.0 # Pin the version you want
```

Or you can download pre-compiled binary from the [releases page](https://github.com/thomasgouveia/vault-aws-credentials-helper/releases).

## Usage

For the following section, we assume that you have a working Vault instance on http://localhost:8200 and that the instance is correctly configured to issue AWS credentials using a Vault role called `developer` and with the authentication method [userpass](https://developer.hashicorp.com/vault/docs/auth/userpass) enabled. 

**You can retrieve the commands required to deploy a such configuration locally in the [examples](./examples/vault-with-userpass/) folder.**

## Without AWS CLI

> The `vault-aws-credentials-helper` can be used without the AWS CLI **for development and test purposes**. For production use, [you must configure your AWS CLI to use it](#with-aws-cli).

Run the following command to test the issuance of credentials with different parameters:

```bash
vault-aws-credentials-helper --vault.addr http://localhost:8200 --vault.auth-method userpass --userpass.username john.smith --userpass.password mysuperpassword --aws.role developer --aws.ttl 30m
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

See next section or use the flag `--help` for more information about flags used.

## With AWS CLI

Create a new profile in your `~/.aws/config` to indicate the AWS CLI to use an external process to retrieve credentials :


```bash
# ~/.aws/config

# ...

[profile developer]
credential_process = /path/to/binary/vault-aws-credentials-helper --vault.addr http://localhost:8200 --vault.auth-method userpass --userpass.username john.smith --userpass.password mysuperpassword --aws.role developer --aws.ttl 30m

# ...
```

In the above command, we use different flags:

- `--vault.addr`: The Vault address to connect.
- `--vault.auth-method`: The authentication method to use to login to Vault.
- `--userpass.username`: The username of the user to authenticate.
- `--userpass.password`: The password of the user to authenticate.
- `--aws.role`: The name of the AWS backend role to use to generate credentials.
- `--aws.ttl`: The TTL of the credentials.

If you have used the configuration of the example given, the `developer` role allows you to access EC2 only in region `eu-west-3`. Let's test this by executing the AWS CLI command:

```bash
aws ec2 describe-instances --region eu-west-3 --profile developer
# Should work! Output redacted.
```

If we update the region to `eu-west-1`, we are not allowed to access EC2:

```bash
aws ec2 describe-instances --region eu-west-1 --profile developer
# An error occurred (UnauthorizedOperation) when calling the DescribeInstances operation: You are not authorized to perform this operation.
```

## License

This project is [MIT licensed](./LICENSE).
