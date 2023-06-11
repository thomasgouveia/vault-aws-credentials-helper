package credentials

import (
	"context"

	"github.com/hashicorp/vault-client-go"
)

type tokenStrategy struct {
	Token string `flag.name:"token" flag.desc:"The token to use to authenticate with Vault. Do not use this authentication method in production."`
}

// Ensure the implementation satisfies the interface
var (
	_ vaultLoginStrategy = &tokenStrategy{}
)

// login performs the login using the token authentication method on Vault.
func (s *tokenStrategy) login(ctx context.Context, client *vault.Client) (*vault.ResponseAuth, error) {
	// As token is a special authentication method, we don't need to perform
	// any authentication process so we can return directly the token wrapped into
	// a vault.ResponseAuth.
	//
	// This authentication method is only available for development and
	// testing purposes and MUST not be used in production.
	return &vault.ResponseAuth{
		ClientToken: s.Token,
	}, nil
}
