package vaultauth

import (
	"context"
	"errors"

	"github.com/hashicorp/vault-client-go"
	"github.com/spf13/cobra"
	"github.com/thomasgouveia/vault-aws-credentials-helper/pkg/structtoflags"
)

// AuthMethod is a type alias for string
// used to define constants for available authentication methods.
type AuthMethod string

const (
	Token    AuthMethod = "token"
	Userpass AuthMethod = "userpass"
	AppRole  AuthMethod = "approle"
)

// vaultLoginStrategy defines the common interface between all the authentication
// methods of Vault supported by the credentials helper.
type vaultLoginStrategy interface {
	// login performs the authentication against Vault and must return a vault.ResponseAuth
	// object or an error if the authentication fails.
	login(ctx context.Context, client *vault.Client) (*vault.ResponseAuth, error)
}

// strategies maps the strategy implementation with the name of the
// authentication method to be used.
var strategies = map[AuthMethod]vaultLoginStrategy{
	Token:    &tokenStrategy{},
	Userpass: &userpassStrategy{},
	AppRole:  &approleStrategy{},
}

var (
	ErrAuthenticationMethodNotSupported = errors.New("vault/auth: authentication method is not supported")
)

// Login is a convenience helper to perform authentication against Vault.
// It handles the validation of the authentication method based on the method
// given in parameters. If the authentication method is not supported yet, an
// error ErrAuthenticationMethodNotSupported will be returned.
func Login(cmd *cobra.Command, client *vault.Client, method AuthMethod) (*vault.ResponseAuth, error) {
	auth, ok := strategies[method]
	if !ok {
		return nil, ErrAuthenticationMethodNotSupported
	}

	// Extract the flags from the command into the auth strategy implementation configuration
	if err := structtoflags.MapCommandFlagsToStruct(cmd, &auth); err != nil {
		return nil, err
	}

	return auth.login(cmd.Context(), client)
}

// MapAuthMethodsConfigToCommandFlags will automatically convert and attach the authentication method
// configurations into the given cobra.Command.
func MapAuthMethodsConfigToCommandFlags(cmd *cobra.Command) error {
	for _, elem := range strategies {
		// Map the strategy flags to the given cobra.Command
		// If an error occurs, we should return it immediately
		if err := structtoflags.MapStructToCommandFlags(cmd, &elem); err != nil {
			return err
		}
	}
	return nil
}
