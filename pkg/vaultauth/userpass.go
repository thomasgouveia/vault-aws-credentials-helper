package vaultauth

import (
	"context"
	"errors"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var (
	ErrUsernameMustNotBeEmpty = errors.New("vault/auth/userpass: username must not be empty")
	ErrPasswordMustNotBeEmpty = errors.New("vault/auth/userpass: password must not be empty")
)

// userpassStrategy defines the configuration that should be attached
// to the command in order to configure this authentication method.
type userpassStrategy struct {
	MountPath string `flag.name:"userpass.mount-path" flag.default:"userpass" flag.desc:"The path to the userpass authentication method in your Vault."`
	Username  string `flag.name:"userpass.username" flag.desc:"The username used to perform login."`
	Password  string `flag.name:"userpass.password" flag.desc:"The password of the user to perform login."`
}

// Ensure the implementation satisfies the interface.
var (
	_ vaultLoginStrategy = &userpassStrategy{}
)

// login performs the login using the userpass authentication method on Vault.
func (s *userpassStrategy) login(ctx context.Context, client *vault.Client) (*vault.ResponseAuth, error) {
	if s.Username == "" {
		return nil, ErrUsernameMustNotBeEmpty
	}

	if s.Password == "" {
		return nil, ErrPasswordMustNotBeEmpty
	}

	opts := []vault.RequestOption{vault.WithMountPath(s.MountPath)}
	req := schema.UserpassLoginRequest{
		Password: s.Password,
	}

	resp, err := client.Auth.UserpassLogin(ctx, s.Username, req, opts...)
	if err != nil {
		return nil, err
	}

	return resp.Auth, nil
}
