package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// err := genKey()
	err := genToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func genToken() error {

	file, err := os.Open("zarf/keys/7e1293da-733d-42f0-9ff5-b2c505c50bdc.pem")
	if err != nil {
		return err
	}

	// limit PEM file size to 1 megabyte. This should be reasonable for almost any PEM file and prevents
	// issues such as linking to other files.
	privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return fmt.Errorf("reading auth private keys: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing auth private keys: %w", err)
	}

	// =========================================================================================================

	// Generating a token requires defining a set of claims. In this applications case, we only care about
	// defining the subject and the user in question and the roles they have on the database.
	// This token will expire in a year
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expirtes
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique dentifier; can be used to prevent the JWT from being replayed (allows token to be used only once)
	claims := struct {
		jwt.StandardClaims
		Roles []string
	}{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   "test-user-id",
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
		Roles: []string{"ADMIN"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// set Key Identifier as we can have multiple keys in rotation and need to know which one to sign.
	// kid tells us what public key to use for signing
	// use the private key file name as the kid
	token.Header["kid"] = "7e1293da-733d-42f0-9ff5-b2c505c50bdc"

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("============ TOKEN BEGIN ============")
	fmt.Println(tokenStr)
	fmt.Println("============ TOKEN END ============")
	fmt.Print("\n")

	// =========================================================================================================

	// Marshal the the public key from the private key to PKIX
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// Construct the PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the private key file.
	if err := pem.Encode(os.Stdout, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	// =========================================================================================================

	// Create the token parser to use. The algorithm used to sign the JWT must be validated to avoid
	// crititical vulnerability.
	parser := jwt.Parser{
		ValidMethods: []string{"RS256"},
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing key id (kid) in token header")
		}
		kidID, ok := kid.(string)
		if !ok {
			return nil, errors.New("user token key id (kid) must be a string")
		}
		fmt.Println("KID: ", kidID)
		return &privateKey.PublicKey, nil
	}

	var parsedClaims struct {
		jwt.StandardClaims
		Roles []string
	}

	parsedToken, err := parser.ParseWithClaims(tokenStr, &parsedClaims, keyFunc)
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	if !parsedToken.Valid {
		return errors.New("invalid token")
	}

	fmt.Println("=============================")
	fmt.Println("Token validated")

	return nil
}

// genKey creates an x509 private/public key for auth tokens
func genKey() error {

	// Private Key
	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Create a file for the private key information in PEM form.
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("creating private file: %w", err)
	}
	defer privateFile.Close()

	// Construct a PEM block for the private key.
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the private key file.
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding to private file: %w", err)
	}

	// =========================================================================================================
	// Public Key

	// Marshal the the public key from the private key to PKIX
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// Create a file for the public key information in PEM form.
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("creating public file: %w", err)
	}
	defer publicFile.Close()

	// Construct the PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the private key file.
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	fmt.Println("private and public key files generated")
	return nil
}
