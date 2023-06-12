package resolver

import (
	"github.com/hashicorp/vault-client-go"
	"github.com/spf13/cobra"
	"github.com/thomasgouveia/vault-aws-credentials-helper/pkg/vaultauth"
)

// resolveConfig holds all the configuration needed to resolve a set of credentials.
type resolveConfig struct {
	// Client is the vault.Client to use to communicate with Vault.
	Client *vault.Client

	// Command is the cobra.Command that holds all the flags necessary to extract
	// configuration and map them into the authentication methods implementation.
	Command *cobra.Command

	// AuthMethod is the authentication method to be used to authenticate against Vault
	//
	// Default: token
	AuthMethod vaultauth.AuthMethod

	// MountPath is the path to the AWS secret engine used to issue credentials in Vault.
	// Default: aws
	MountPath string

	// Role is the name of the Vault role to use in the AWS backend to issue credentials.
	Role string

	// TTL is the time-to-leave of the generated credentials.
	TTL string
}

type ResolveOption func(c *resolveConfig) error

// newResolveConfig instantiates a new ResolveConfig with the default configuration
// and override fields with the given options if needed.
func newResolveConfig(opts ...ResolveOption) (*resolveConfig, error) {
	c := &resolveConfig{
		AuthMethod: vaultauth.Token,
		MountPath:  "aws",
		TTL:        "15m",
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// WithCommand specifies the cobra.Command that calls the resolver.
func WithCommand(cmd *cobra.Command) ResolveOption {
	return func(c *resolveConfig) error {
		c.Command = cmd
		return nil
	}
}

// WithClient specifies the vault.Client to use for communicating with Vault.
func WithClient(client *vault.Client) ResolveOption {
	return func(c *resolveConfig) error {
		c.Client = client
		return nil
	}
}

// WithAuthMethod specifies the authentication method to use
// to authenticate with Vault.
//
// Default: token.
func WithAuthMethod(method vaultauth.AuthMethod) ResolveOption {
	return func(c *resolveConfig) error {
		c.AuthMethod = method
		return nil
	}
}

// WithMountPath specifies the path of the AWS backend to use in Vault.
//
// Default: aws.
func WithMountPath(mountPath string) ResolveOption {
	return func(c *resolveConfig) error {
		c.MountPath = mountPath
		return nil
	}
}

// WithRole specifies the role to use in the AWS backend to issue credentials.
func WithRole(role string) ResolveOption {
	return func(c *resolveConfig) error {
		c.Role = role
		return nil
	}
}

// WithTTL specifies the TTL for the generated credentials.
//
// Default: 15m.
func WithTTL(ttl string) ResolveOption {
	return func(c *resolveConfig) error {
		c.TTL = ttl
		return nil
	}
}
