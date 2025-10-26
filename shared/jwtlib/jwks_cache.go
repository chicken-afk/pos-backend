package jwtlib

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var (
	jwksCache JWKS
	lastFetch time.Time
	mu        sync.RWMutex
)

func GetJWKS(url string) (JWKS, error) {
	mu.RLock()
	if time.Since(lastFetch) < 15*time.Minute && len(jwksCache.Keys) > 0 {
		defer mu.RUnlock()
		return jwksCache, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()
	resp, err := http.Get(url)
	if err != nil {
		return JWKS{}, err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JWKS{}, err
	}
	jwksCache = jwks
	lastFetch = time.Now()
	return jwks, nil
}
