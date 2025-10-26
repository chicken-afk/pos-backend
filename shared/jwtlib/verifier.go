package jwtlib

import (
	"encoding/pem"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenString, jwksURL string) (*jwt.Token, error) {
	jwks, err := GetJWKS(jwksURL)
	if err != nil {
		return nil, err
	}

	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid")
		}
		for _, key := range jwks.Keys {
			if key.Kid == kid {
				block, _ := pem.Decode([]byte(key.Pem))
				pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pem.EncodeToMemory(block))
				if err != nil {
					return nil, err
				}
				return pubKey, nil
			}
		}
		return nil, errors.New("no matching key")
	})
}
