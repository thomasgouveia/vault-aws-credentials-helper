package resolver

import (
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/thomasgouveia/vault-aws-credentials-helper/pkg/awscreds"
	"github.com/thomasgouveia/vault-aws-credentials-helper/pkg/vaultauth"
)

var (
	ErrVaultRoleEmpty = errors.New("the vault role must not be empty")
)

// FetchCredentials authenticate with Vault using the given authentication method, and ask for AWS credentials
// on the underlying AWS secret engine.
func ResolveCredentials(opts ...ResolveOption) (*awscreds.AWSCredentials, error) {
	cfg, err := newResolveConfig(opts...)
	if err != nil {
		return nil, err
	}

	ctx := cfg.Command.Context()

	// Ensure the user has given a Vault role
	if cfg.Role == "" {
		return nil, ErrVaultRoleEmpty
	}

	// Perform the authentication with Vault
	auth, err := vaultauth.Login(cfg.Command, cfg.Client, cfg.AuthMethod)
	if err != nil {
		return nil, err
	}

	// Ask the AWS secret engine to fetch credentials
	// Currently, only AWS STS based credentials are supported
	reqOpts := []vault.RequestOption{
		vault.WithToken(auth.ClientToken),
		vault.WithMountPath(cfg.MountPath),
		vault.WithCustomQueryParameters(url.Values{
			"ttl": {cfg.TTL},
		}),
	}

	resp, err := cfg.Client.Secrets.AwsGenerateStsCredentials(ctx, cfg.Role, reqOpts...)
	if err != nil {
		return nil, err
	}

	return mapVaultResponseToAWSCredentials(resp.Data)
}

// mapVaultResponseToAWSCredentials is a helper function to map a Vault API response
// into a pair of awscreds.AWSCredentials.
func mapVaultResponseToAWSCredentials(data map[string]interface{}) (*awscreds.AWSCredentials, error) {
	opts := []awscreds.AWSCredentialOption{}

	if accessKey, ok := data["access_key"].(string); ok {
		opts = append(opts, awscreds.WithAccessKeyId(accessKey))
	}

	if secretKey, ok := data["secret_key"].(string); ok {
		opts = append(opts, awscreds.WithSecretAccessKey(secretKey))
	}

	if sessionToken, ok := data["security_token"].(string); ok {
		opts = append(opts, awscreds.WithSessionToken(sessionToken))
	}

	if ttl, ok := data["ttl"].(json.Number); ok {
		ttl, err := ttl.Int64()
		if err != nil {
			return nil, err
		}

		expiry := time.Now().Add(time.Duration(ttl) * time.Second)
		opts = append(opts, awscreds.WithExpiration(expiry))
	}

	return awscreds.New(opts...)
}
