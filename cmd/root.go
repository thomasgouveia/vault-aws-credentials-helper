package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/spf13/cobra"
	"github.com/thomasgouveia/vault-aws-credentials-helper/pkg/resolver"
	"github.com/thomasgouveia/vault-aws-credentials-helper/pkg/vaultauth"
)

var rootCmd = &cobra.Command{
	Use:   "vault-aws-credentials-helper",
	Short: "Configure your AWS CLI to retrieve dynamic and short-lived AWS credentials from HashiCorp Vault with ease.",
	Long:  "Configure your AWS CLI to retrieve dynamic and short-lived AWS credentials from HashiCorp Vault with ease.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Retrieve global flags value from command
		vaultAddr, _ := cmd.Flags().GetString("vault.addr")
		vaultAuthMethod, _ := cmd.Flags().GetString("vault.auth-method")

		// AWS credentials config
		awsMountPath, _ := cmd.Flags().GetString("aws.mount-path")
		awsRole, _ := cmd.Flags().GetString("aws.role")
		awsTtl, _ := cmd.Flags().GetString("aws.ttl")

		client, err := vault.New(vault.WithAddress(vaultAddr), vault.WithRequestTimeout(30*time.Second))
		if err != nil {
			return err
		}

		// Define options for resolving credentials
		resolveOpts := []resolver.ResolveOption{
			resolver.WithCommand(cmd),
			resolver.WithClient(client),
			resolver.WithAuthMethod(vaultauth.AuthMethod(vaultAuthMethod)),
			resolver.WithMountPath(awsMountPath),
			resolver.WithRole(awsRole),
			resolver.WithTTL(awsTtl),
		}

		// Login to Vault, and then issue the credentials
		// by calling the correct endpoint depending on
		// the given options.
		creds, err := resolver.ResolveCredentials(resolveOpts...)
		if err != nil {
			if errors.Is(err, resolver.ErrVaultRoleEmpty) {
				return fmt.Errorf("you must provide a Vault role configured in your AWS backend to generate credentials using --aws.role")
			}
			return err
		}

		out, err := creds.JSONString()
		if err != nil {
			return err
		}

		fmt.Println(out)
		return nil
	},
}

func init() {
	// Vault global flags
	rootCmd.PersistentFlags().String("vault.addr", "https://127.0.0.1:8200", "The address of the Vault where to perform login and credentials generation.")
	rootCmd.PersistentFlags().String("vault.auth-method", "userpass", "The authentication method to use to authenticate with Vault.")

	// Credentials configuration
	rootCmd.PersistentFlags().String("aws.mount-path", "aws", "The mount path of the AWS backend engine to use.")
	rootCmd.PersistentFlags().String("aws.role", "", "The name of the Vault role to use to generate credentials on the AWS backend.")
	rootCmd.PersistentFlags().String("aws.ttl", "15m", "The TTL of the Vault lease for the AWS generated credentials.")

	// Bind authentication methods flags
	if err := vaultauth.MapAuthMethodsConfigToCommandFlags(rootCmd); err != nil {
		panic(err)
	}
}

// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
