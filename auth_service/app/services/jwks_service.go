package services

import "log"

type JwksService interface {
	GetJwks() (JWKS, error)
}

type jwksService struct {
	publicDATAJWT []byte
}

func NewJwksService(publicDATAJWT []byte) JwksService {
	return &jwksService{
		publicDATAJWT: publicDATAJWT,
	}
}

type JWKS struct {
	Keys []map[string]string `json:"keys"`
}

func (s *jwksService) GetJwks() (JWKS, error) {

	log.Println("Retrieving JWKS")
	log.Println(s.publicDATAJWT)

	// Convert public PEM to base64 (simplified JWKS format)
	jwks := JWKS{
		Keys: []map[string]string{
			{
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"kid": "rsa-key-1",
				"pem": string(s.publicDATAJWT),
			},
		},
	}
	return jwks, nil
}
