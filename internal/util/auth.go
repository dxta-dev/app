package util

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
)

type Key struct {
	Crv string `json:"crv"`
	X   string `json:"x"`
	Kty string `json:"kty"`
	Kid string `json:"kid"`
}

type VerificationKeys struct {
	Keys []Key `json:"keys"`
}

func CreateAuthVerifier(pubKey ed25519.PublicKey) *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("EdDSA", nil, pubKey)
	return tokenAuth
}

func GetRawPublicKeyFromJWK(jwkData Key) (ed25519.PublicKey, error) {
	decodedKey, err := base64.RawURLEncoding.DecodeString(jwkData.X)
	if err != nil {
		return nil, fmt.Errorf("failed to base64url decode public key material: %w", err)
	}

	if len(decodedKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid Ed25519 public key length: got %d, want %d", len(decodedKey), ed25519.PublicKeySize)
	}

	return ed25519.PublicKey(decodedKey), nil
}

func RetrieveJWK() (Key, error) {
	resp, err := http.Get(os.Getenv("AUTH_JWKS_URL"))

	if err != nil {
		return Key{}, err
	}

	if resp.StatusCode != 200 {
		return Key{}, fmt.Errorf("failed to retrieve jwks. Status: %s", resp.Status)
	}

	defer resp.Body.Close()

	j := VerificationKeys{}
	err = json.NewDecoder(resp.Body).Decode(&j)

	if err != nil {
		return Key{}, err
	}

	if len(j.Keys) == 0 {
		return Key{}, fmt.Errorf("no keys found in JWKS")
	}

	return j.Keys[0], nil
}
