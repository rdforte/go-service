// Package keystore implements the auth.Keystore interface. This implements
// an in-memory keystore for JWT support
package keystore

import (
	"crypto/rsa"
	"errors"
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
type KeyStore struct {
	mu    sync.RWMutex
	store map[string]*rsa.PrivateKey
}

// New constructs an empty KeyStore ready for use.
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

// NewMap constructs a Keystore with an initial set of keys.
func NewMap(store map[string]*rsa.PrivateKey) *KeyStore {
	return &KeyStore{
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
func NewFS() (*KeyStore, error) {
	ks := KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}

	// For Debug purposes read the local secret
	// Note that when running in docker the current WORKDIR is app/services/sales-api
	file, err := os.Open("../../../zarf/keys/7e1293da-733d-42f0-9ff5-b2c505c50bdc.pem")
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

	ks.store["7e1293da-733d-42f0-9ff5-b2c505c50bdc"] = privateKey

	return &ks, nil

}

// Add adds a private key and combination kid to the store.
func (ks *KeyStore) Add(privateKey *rsa.PrivateKey, kid string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.store[kid] = privateKey
}

// Remove removes a private key and combination kid to the store.
func (ks *KeyStore) Remove(kid string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	delete(ks.store, kid)
}

// PrivateKey searches the key store for a given kid and returns
// the private key.
func (ks *KeyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, errors.New("kid lookup failed")
	}
	return privateKey, nil
}

// PublicKey searches the key store for a given kid and returns
// the public key.
func (ks *KeyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, errors.New("kid lookup failed")
	}
	return &privateKey.PublicKey, nil
}
