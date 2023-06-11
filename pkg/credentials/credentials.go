package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/spf13/cobra"
	"github.com/thomasgouveia/vault-aws-credentials-helper/pkg/structtoflags"
)

// AWSCredentials represents the JSON structure that the
// vault-aws-credentials-helper should output to STDOUT to allow
// the AWS CLI to autoconfigure with the generated credentials.
//
// See https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html
type AWSCredentials struct {
	// Version must be set to 1 by default. This might increment over
	// time if the AWS credentials structure evolves.
	Version int `json:"Version"`
	// AccessKeyId the access key generated.
	AccessKeyId string `json:"AccessKeyId"`
	// SecretAccessKey the secret access key generated.
	SecretAccessKey string `json:"SecretAccessKey"`
	// SessionToken the AWS session token for temporary credentials
	SessionToken string `json:"SessionToken"`
	// Expiration should be an ISO8601 formatted timestamp.
	// If the Expiration key is not present in the output, the AWS CLI
	// will assume that the credentials are long-term credentials that do not refresh.
	// Otherwise the credentials are considered temporary credentials and are refreshed
	// automatically by rerunning the credential_process command before they expire.
	Expiration string `json:"Expiration,omitempty"`
}

// vaultLoginStrategy defines the common interface between all the authentication
// methods of Vault supported by the credentials helper.
type vaultLoginStrategy interface {
	// login performs the authentication against Vault and must return a vault.ResponseAuth
	// object or an error if the authentication fails.
	login(ctx context.Context, client *vault.Client) (*vault.ResponseAuth, error)
}

// authStrategies maps the strategy implementation with the name of the
// authentication method to be used.
var authStrategies = map[string]vaultLoginStrategy{
	"token":    &tokenStrategy{},
	"userpass": &userpassStrategy{},
}

var (
	ErrUnknownAuthMethod = errors.New("unknown auth method")
)

type FetchCredentialConfig struct {
	// AuthMethod is the authentication method to be used to authenticate against Vault
	AuthMethod string
	// MountPath is the path to the AWS secret engine used to issue credentials
	MountPath string
	// Role is the name of the role allowed to issue credentials on this backend.
	Role string
	// TTL is the duration of the generated credentials
	TTL string
}

// Fetch authenticate with Vault using the given authentication method, and ask for AWS credentials
// on the underlying AWS secret engine.
func Fetch(cmd *cobra.Command, client *vault.Client, cfg *FetchCredentialConfig) (*AWSCredentials, error) {
	ctx := cmd.Context()

	strategy, ok := authStrategies[cfg.AuthMethod]
	if !ok {
		return nil, ErrUnknownAuthMethod
	}

	// Extract the flags from the command into the strategy implementation configuration
	if err := structtoflags.MapCommandFlagsToStruct(cmd, &strategy); err != nil {
		return nil, err
	}

	// Perform the authentication with Vault
	auth, err := strategy.login(ctx, client)
	if err != nil {
		return nil, err
	}

	// Ask the AWS secret engine to fetch credentials
	// Currently, only AWS STS based credentials are supported
	opts := []vault.RequestOption{
		vault.WithToken(auth.ClientToken),
		vault.WithMountPath(cfg.MountPath),
		vault.WithCustomQueryParameters(url.Values{
			"ttl": {cfg.TTL},
		}),
	}

	resp, err := client.Secrets.AwsGenerateStsCredentials(ctx, cfg.Role, opts...)
	if err != nil {
		return nil, err
	}

	ttl, _ := resp.Data["ttl"].(json.Number).Int64()
	expiry := time.Now().Add(time.Duration(ttl) * time.Second)

	return &AWSCredentials{
		Version:         1,
		AccessKeyId:     resp.Data["access_key"].(string),
		SecretAccessKey: resp.Data["secret_key"].(string),
		SessionToken:    resp.Data["security_token"].(string),
		Expiration:      expiry.Format("2006-01-02T15:04:05.000Z"),
	}, nil
}

// ConfigureAuthFlags
func ConfigureAuthFlags(cmd *cobra.Command) error {
	for _, elem := range authStrategies {
		// Map the strategy flags to the given cobra.Command
		// If an error occurs, we should return it immediately
		if err := structtoflags.MapStructToCommandFlags(cmd, &elem); err != nil {
			return err
		}
	}
	return nil
}
