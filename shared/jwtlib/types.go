package jwtlib

import "github.com/golang-jwt/jwt/v5"

// JWK merepresentasikan satu key publik di JWKS (public endpoint)
type JWK struct {
	Kty string `json:"kty"` // Key type (ex: RSA)
	Alg string `json:"alg"` // Algorithm (ex: RS256)
	Use string `json:"use"` // Usage (ex: sig)
	Kid string `json:"kid"` // Key ID (unique identifier)
	Pem string `json:"pem"` // Public key dalam format PEM
}

// JWKS adalah kumpulan JWK
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// CustomClaims contoh jika nanti kamu ingin extract data tertentu dari JWT
type CustomClaims struct {
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
}
