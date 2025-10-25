// Package auth provides JWT token generation and validation.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// GoogleTokenInfo represents Google's token info response
type GoogleTokenInfo struct {
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"` // Google user ID
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Exp           string `json:"exp"`
}

// GoogleAuthVerifier verifies Google Sign In tokens
type GoogleAuthVerifier struct {
	httpClient *http.Client
}

const (
	// #nosec G101 -- This is a public Google API endpoint, not a credential
	googleTokenInfoURL = "https://oauth2.googleapis.com/tokeninfo"
)

// NewGoogleAuthVerifier creates a new Google auth verifier
func NewGoogleAuthVerifier() *GoogleAuthVerifier {
	return &GoogleAuthVerifier{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyIDToken verifies a Google ID token and returns the token info
func (v *GoogleAuthVerifier) VerifyIDToken(idToken string) (*GoogleTokenInfo, error) {
	// Build request URL with ID token
	url := fmt.Sprintf("%s?id_token=%s", googleTokenInfoURL, idToken)

	// Create request with context
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token verification failed with status: %d", resp.StatusCode)
	}

	// Decode response
	var tokenInfo GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %w", err)
	}

	// Verify email is verified
	if tokenInfo.EmailVerified != "true" {
		return nil, errors.New("email not verified")
	}

	return &tokenInfo, nil
}
