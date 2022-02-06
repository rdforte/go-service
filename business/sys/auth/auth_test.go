// Package auth_test
// We call the package auth_test instead of auth because then it forces us to import our auth package
// and implement it how any other user would.
package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rdforte/go-service/business/sys/auth"
)

// Success and failure markers.
// Helped to increase verbosity in our tests.
const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestAuth(t *testing.T) {
	t.Log("Given the need to be able to authenticate and authorize access.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a singe user.", testID)
		{
			// Setup Private Key
			const keyID = "8e1293db-733d-42f0-9ff9-b2c505c50bdc"
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a private key: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a private key.", success, testID)

			// Construct Auth
			a, err := auth.New(keyID, &keyStore{pk: privateKey})
			if err != nil {
				t.Errorf("\t%s\tTest %d:\tShould be able to create an authenticator: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create an authenticator.", success, testID)

			// Define Claims
			claims := auth.Claims{
				StandardClaims: jwt.StandardClaims{
					Issuer:    "service project",
					Subject:   "00000000-0000-0000-0000-000000000000",
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
					IssuedAt:  time.Now().UTC().Unix(),
				},
				Roles: []string{auth.RoleAdmin},
			}

			// Generate Token
			token, err := a.GenerateToken(claims)
			if err != nil {
				t.Errorf("\t%s\tTest %d:\tShould be able to generate a JWT: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to generate a JWT.", success, testID)

			// Validate Token
			parsedClaims, err := a.ValidateToken(token)
			if err != nil {
				t.Errorf("\t%s\tTest %d:\tShould be able to parse the claims: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to parse the claims.", success, testID)

			// Check the Roles Length are the same
			if exp, got := len(claims.Roles), len(parsedClaims.Roles); exp != got {
				t.Logf("\t\tTest %d:\texp: %d", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %d", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShoud have the expected number of roles: %v", failed, testID, len(claims.Roles))
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected number of roles.", success, testID)

			// Check the Roles are the same
			if exp, got := claims.Roles[0], parsedClaims.Roles[0]; exp != got {
				t.Logf("\t\tTest %d:\texp: %s", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %s", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShoud have the expected roles: %v", failed, testID, len(claims.Roles))
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected roles.", success, testID)
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
