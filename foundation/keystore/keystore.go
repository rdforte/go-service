// Package keystore implements the auth.Keystore interface. This implements
// an in-memory keystore for JWT support
package keystore

import (
	"crypto/rsa"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	if _, err := NewFS(); err != nil {
		fmt.Println(err)
	}
}

// KeyStore represents an in memory store implementation of the
// KeyStorer interface for use with the auth package
type Keystore struct {
	mu    sync.RWMutex
	store map[string]*rsa.PrivateKey
}

// New constructs an empty KeyStore ready for use.
func New() *Keystore {
	return &Keystore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

// NewMap constructs a Keystore with an initial set of keys.
func NewMap(store map[string]*rsa.PrivateKey) *Keystore {
	return &Keystore{
		store: store,
	}
}

// NewFS constructs a KeyStore based on a set of PEM files rooted inside of a directory.
// The name of each PEM file will be used as the key id.
// In Production we would have NewFs talk to AWS Secrects Manager to get our PEM files.
// https://aws.amazon.com/blogs/security/how-to-use-aws-secrets-manager-securely-store-rotate-ssh-key-pairs/
// You would need to make a call to `ListSecrets` to get the secrets then
// To get the secret value from SecretString or SecretBinary, call `GetSecretValue`.
// Once you have the secrets you can then set them up in your keystore to be used by the application.
// If you wanted to rotate the secrets every month or so you could setup a Lambda that will generate
// a New Secret which we will then use to sign our New JWT's.
func NewFS() (*Keystore, error) {
	ks := Keystore{
		store: make(map[string]*rsa.PrivateKey),
	}

	// For Debug purposes read the local secret
	file, err := os.Open("./zarf/keys/local-secret.pem")
	if err != nil {
		return nil, fmt.Errorf("NewFS: error opening file: %w", err)
	}
	defer file.Close()

	// limit PEM file size to 1 megabyte. This should be reasonable for almost any PEM file and prevents
	// linking to random files.
	privatePem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return nil, fmt.Errorf("NewFS: error reading file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		return nil, fmt.Errorf("NewFS: error parsing rsa private key from pem: %w", err)
	}

	ks.store["local-secret"] = privateKey

	return &ks, nil

}
