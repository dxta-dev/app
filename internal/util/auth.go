package util

import (
	"context"
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

type JWKS struct {
	Keys []Key `json:"keys"`
}

func retrieveJWK() (Key, error) {
	resp, err := http.Get(os.Getenv("AUTH_JWKS_URL"))

	if err != nil {
		return Key{}, err
	}

	if resp.StatusCode != 200 {
		return Key{}, fmt.Errorf("failed to retrieve jwks. Status: %s", resp.Status)
	}

	defer resp.Body.Close()

	jwks := JWKS{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return Key{}, err
	}

	if len(jwks.Keys) == 0 {
		return Key{}, fmt.Errorf("no keys found in JWKS")
	}

	return jwks.Keys[0], nil
}

func getRawPublicKeyFromJWK(jwkData Key) (ed25519.PublicKey, error) {
	decodedKey, err := base64.RawURLEncoding.DecodeString(jwkData.X)
	if err != nil {
		return nil, fmt.Errorf("failed to base64url decode public key material: %w", err)
	}

	if len(decodedKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid Ed25519 public key length: got %d, want %d", len(decodedKey), ed25519.PublicKeySize)
	}

	return ed25519.PublicKey(decodedKey), nil
}

func GetRawPublicKey() (ed25519.PublicKey, error) {
	jwkData, err := retrieveJWK()

	if err != nil {
		fmt.Println("Could not retrieve JWK")
		return nil, err
	}

	fmt.Println("Could not get public key from JWK")
	pubKey, err := getRawPublicKeyFromJWK(jwkData)

	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func CreateAuthVerifier(pubKey ed25519.PublicKey) *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("EdDSA", nil, pubKey)
	return tokenAuth
}

type contextKey struct {
	key string
}

var (
	AuthIdCtxKey = contextKey{"authId"}
)

func Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := jwtauth.FromContext(r.Context())

			if err != nil {
				fmt.Printf("Error extracting token and claims from context. Error: %s", err.Error())
				JSONError(w, ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
				return
			}

			if token == nil {
				fmt.Println("No token found in context")
				JSONError(w, ErrorParam{Error: "Unauthorized"}, http.StatusUnauthorized)
				return
			}

			authId := claims["organizationId"]

			if authId == nil {
				fmt.Println("No auth id found in JWT payload")
				JSONError(w, ErrorParam{Error: "Bad request"}, http.StatusBadRequest)
				return
			}

			authId = authId.(string)

			ctx := r.Context()

			ctx = context.WithValue(ctx, AuthIdCtxKey, authId)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
