// Package auth provides JWT token generation and validation.
package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ApplePublicKey represents Apple's public key for JWT verification
type ApplePublicKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// ApplePublicKeys represents the response from Apple's public keys endpoint
type ApplePublicKeys struct {
	Keys []ApplePublicKey `json:"keys"`
}

// AppleIdentityToken represents the claims in Apple's identity token
type AppleIdentityToken struct {
	Iss            string `json:"iss"`
	Aud            string `json:"aud"`
	Exp            int64  `json:"exp"`
	Iat            int64  `json:"iat"`
	Sub            string `json:"sub"` // Apple user ID
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"`
	IsPrivateEmail string `json:"is_private_email"`
	jwt.RegisteredClaims
}

// AppleAuthVerifier verifies Apple Sign In identity tokens
type AppleAuthVerifier struct {
	publicKeys  map[string]*rsa.PublicKey
	lastFetched time.Time
	httpClient  *http.Client
}

const (
	applePublicKeysURL = "https://appleid.apple.com/auth/keys"
	appleIssuer        = "https://appleid.apple.com"
	keyRefreshInterval = 24 * time.Hour
)

// NewAppleAuthVerifier creates a new Apple auth verifier
func NewAppleAuthVerifier() *AppleAuthVerifier {
	return &AppleAuthVerifier{
		publicKeys: make(map[string]*rsa.PublicKey),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyIdentityToken verifies an Apple identity token and returns the claims
func (v *AppleAuthVerifier) VerifyIdentityToken(identityToken string, clientID string) (*AppleIdentityToken, error) {
	// Refresh keys if needed
	if time.Since(v.lastFetched) > keyRefreshInterval || len(v.publicKeys) == 0 {
		if err := v.fetchPublicKeys(); err != nil {
			return nil, fmt.Errorf("failed to fetch Apple public keys: %w", err)
		}
	}

	// Parse the token to get the key ID (kid) from header
	token, err := jwt.ParseWithClaims(identityToken, &AppleIdentityToken{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the key ID from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found in token header")
		}

		// Get the public key for this kid
		publicKey, ok := v.publicKeys[kid]
		if !ok {
			return nil, fmt.Errorf("public key not found for kid: %s", kid)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate token
	claims, ok := token.Claims.(*AppleIdentityToken)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Verify issuer
	if claims.Iss != appleIssuer {
		return nil, fmt.Errorf("invalid issuer: %s", claims.Iss)
	}

	// Verify audience (client ID)
	if claims.Aud != clientID {
		return nil, fmt.Errorf("invalid audience: %s", claims.Aud)
	}

	// Verify expiration
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// fetchPublicKeys fetches Apple's public keys for JWT verification
func (v *AppleAuthVerifier) fetchPublicKeys() error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, applePublicKeysURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch public keys: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var keys ApplePublicKeys
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return fmt.Errorf("failed to decode public keys: %w", err)
	}

	// Convert JWK to RSA public keys
	newKeys := make(map[string]*rsa.PublicKey)
	for _, key := range keys.Keys {
		publicKey, err := v.jwkToRSAPublicKey(key)
		if err != nil {
			return fmt.Errorf("failed to convert JWK to RSA public key: %w", err)
		}
		newKeys[key.Kid] = publicKey
	}

	v.publicKeys = newKeys
	v.lastFetched = time.Now()

	return nil
}

// jwkToRSAPublicKey converts a JWK to an RSA public key
func (v *AppleAuthVerifier) jwkToRSAPublicKey(key ApplePublicKey) (*rsa.PublicKey, error) {
	// Decode the modulus (n)
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	// Decode the exponent (e)
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert bytes to big integers
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// Create RSA public key
	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}
