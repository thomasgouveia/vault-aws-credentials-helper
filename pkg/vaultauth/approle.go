package vaultauth

import (
	"context"
	"errors"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var (
	ErrRoleIdMustNotBeEmpty   = errors.New("vault/auth/approle: the role-id must not be empty")
	ErrSecretIdMustNotBeEmpty = errors.New("vault/auth/approle: the secret-id must not be empty")
)

// approleStrategy defines the configuration that should be attached
// to the command in order to configure the AppRole authentication.
type approleStrategy struct {
	MountPath string `flag.name:"approle.mount-path" flag.default:"approle" flag.desc:"The path to the AppRole authentication method in your Vault."`
	RoleId    string `flag.name:"approle.role-id" flag.desc:"The identifier of the role to use to perform the login."`
	SecretId  string `flag.name:"approle.secret-id" flag.desc:"The secret identifier of the role to use to perform the login."`
}

// Ensure the implementation satisfies the interface.
var (
	_ vaultLoginStrategy = &approleStrategy{}
)

// login performs the login using the AppRole authentication method on Vault.
func (s *approleStrategy) login(ctx context.Context, client *vault.Client) (*vault.ResponseAuth, error) {
	if s.RoleId == "" {
		return nil, ErrRoleIdMustNotBeEmpty
	}

	if s.SecretId == "" {
		return nil, ErrSecretIdMustNotBeEmpty
	}

	opts := []vault.RequestOption{vault.WithMountPath(s.MountPath)}
	req := schema.AppRoleLoginRequest{
		RoleId:   s.RoleId,
		SecretId: s.SecretId,
	}

	resp, err := client.Auth.AppRoleLogin(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	return resp.Auth, nil
}
