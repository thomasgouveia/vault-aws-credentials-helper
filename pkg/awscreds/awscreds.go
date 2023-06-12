package awscreds

import (
	"encoding/json"
	"time"
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
	AccessKeyId string `json:"AccessKeyId,omitempty"`
	// SecretAccessKey the secret access key generated.
	SecretAccessKey string `json:"SecretAccessKey,omitempty"`
	// SessionToken the AWS session token for temporary credentials
	SessionToken string `json:"SessionToken,omitempty"`
	// Expiration should be an ISO8601 formatted timestamp.
	// If the Expiration key is not present in the output, the AWS CLI
	// will assume that the credentials are long-term credentials that do not refresh.
	// Otherwise the credentials are considered temporary credentials and are refreshed
	// automatically by rerunning the credential_process command before they expire.
	Expiration string `json:"Expiration,omitempty"`
}

type AWSCredentialOption func(c *AWSCredentials) error

// New initializes a new AWS credentials with the default configuration.
func New(opts ...AWSCredentialOption) (*AWSCredentials, error) {
	creds := &AWSCredentials{Version: 1}

	// Apply all options
	for _, opt := range opts {
		if err := opt(creds); err != nil {
			return nil, err
		}
	}

	return creds, nil
}

// JSONString returns a pretty JSON string representation of the credentials.
func (c *AWSCredentials) JSONString() (string, error) {
	by, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return "", err
	}
	return string(by), nil
}

// WithAccessKeyId specifies the AWS access key id for this pair of credentials.
func WithAccessKeyId(accessKeyId string) AWSCredentialOption {
	return func(c *AWSCredentials) error {
		c.AccessKeyId = accessKeyId
		return nil
	}
}

// WithSecretAccessKey specifies the AWS secret access key id for this pair of credentials.
func WithSecretAccessKey(secretAccessKey string) AWSCredentialOption {
	return func(c *AWSCredentials) error {
		c.SecretAccessKey = secretAccessKey
		return nil
	}
}

// WithSessionToken specifies the AWS session token this pair of credentials.
func WithSessionToken(sessionToken string) AWSCredentialOption {
	return func(c *AWSCredentials) error {
		c.SessionToken = sessionToken
		return nil
	}
}

// WithExpiration sets the expiration of this pair of credentials.
func WithExpiration(exp time.Time) AWSCredentialOption {
	return func(c *AWSCredentials) error {
		c.Expiration = exp.Format("2006-01-02T15:04:05.000Z")
		return nil
	}
}
