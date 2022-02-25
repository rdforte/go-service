// Package auth_test
// We call the package auth_test instead of auth because then it forces us to import our auth package
// and implement it how any other user would.
package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/foundation/logger"
)

func TestAuth(t *testing.T) {
	tl := logger.NewTestLog(t)

	tl.Describe("Authenticate and authorize access")
	{
		tl.It("should authenticate when handling a single user.")
		{
			// Setup Private Key
			const keyID = "8e1293db-733d-42f0-9ff9-b2c505c50bdc"
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				tl.Failed("Should be able to create a private key", err)
			}
			tl.Success("Should be able to create a private key")

			// Construct Auth
			a, err := auth.New(keyID, &keyStore{pk: privateKey})
			if err != nil {
				tl.Failed("Should be able to create an authenticator", err)
			}
			tl.Success("Should be able to create an authenticator")

			// Define Claims
			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "service project",
					Subject:   "00000000-0000-0000-0000-000000000000",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: []string{auth.RoleAdmin},
			}

			// Generate Token
			token, err := a.GenerateToken(claims)
			if err != nil {
				tl.Failed("Should be able to generate a JWT", err)
			}
			tl.Success("Should be able to generate a JWT")

			// Validate Token
			parsedClaims, err := a.ValidateToken(token)
			if err != nil {
				tl.Failed("Should be able to parse the claims", err)
			}
			tl.Success("Should be able to parse the claims")

			// Check the Roles Length are the same
			if exp, got := len(claims.Roles), len(parsedClaims.Roles); exp != got {
				tl.Failed("Shoud have the expected number of roles", fmt.Errorf("[#roles: %v]", len(claims.Roles)))
			}
			tl.Success("Should have the expected number of roles")

			// Check the Roles are the same
			if exp, got := claims.Roles[0], parsedClaims.Roles[0]; exp != got {
				tl.Failed("Shoud have the expected roles", fmt.Errorf("[roles: %v]", claims.Roles))
			}
			tl.Success("Should have the expected roles")
		}
	}

}

// ===========================================================================================================
// Mocked KeyStore

type keyStore struct {
	pk *rsa.PrivateKey
}

func (ks *keyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	return ks.pk, nil
}

func (ks *keyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	return &ks.pk.PublicKey, nil
}
