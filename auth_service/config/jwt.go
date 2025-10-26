package config

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKeyJWT *jwt.SigningMethodRSA
	keyDataJWT    []byte
	publicDataJWT []byte
)

func InitJWT() {
	// Initialization code for JWT
	priv, err := os.ReadFile("jwks/private.pem")
	if err != nil {
		panic(err)
	}
	pub, err := os.ReadFile("jwks/public.pem")
	if err != nil {
		panic(err)
	}

	keyDataJWT = priv
	publicDataJWT = pub
}
