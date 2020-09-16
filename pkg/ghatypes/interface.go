package ghatypes

import (
	"context"

	"golang.org/x/xerrors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// ErrSecretNotFound represents the specified secret does not exist
	ErrSecretNotFound = xerrors.Errorf("secret is not found")

	// ErrKeyNotFound represents the specified key does not exist in the secret
	ErrKeyNotFound = xerrors.Errorf("key is not found in the secret")
)

// GitHubAppInterface is an interface for GitHub App
type GitHubAppInterface interface {
	// GetURL returns https://github.com or GHE url
	GetURL() string
	// GetAppID returns app id for GitHub App
	GetAppID() int64
	// GetPrivateKey returns the raw private key for GitHub App
	GetPrivateKey(ctx context.Context, c client.Client) ([]byte, error)
}
