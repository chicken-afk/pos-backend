package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type SetRefreshTokenParams struct {
	SessionID    string
	RefreshToken string
	Email        string
	ExpiryTime   int
}

func SetRedisData(redisClient *redis.Client, key string, stringValue string, expiryTime int) error {
	err := redisClient.Set(context.Background(), key, stringValue, time.Duration(expiryTime)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetRefreshToken(redisClient *redis.Client, params SetRefreshTokenParams) error {
	key := "refresh_token:" + params.SessionID + ":" + params.Email
	val := params.RefreshToken
	err := redisClient.Set(context.Background(), key, val, time.Duration(params.ExpiryTime)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func BlacklistAccessToken(redisClient *redis.Client, accessToken string, expiryTime int) error {
	key := "blacklist_token:" + accessToken
	log.Printf("Blacklisting access token in Redis with expiry time: %d seconds\n", expiryTime)
	err := redisClient.Set(context.Background(), key, "blacklisted", time.Duration(expiryTime)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func IsTokenBlacklisted(redisClient *redis.Client, accessToken string) (bool, error) {
	key := "blacklist_token:" + accessToken
	val, err := redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if val == "blacklisted" {
		return true, nil
	}
	return false, nil
}

func RemoveRefreshToken(redisClient *redis.Client, sessionID string, email string) error {
	key := "refresh_token:" + sessionID + ":" + email
	err := redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetRefreshToken(redisClient *redis.Client, sessionID string, email string) (string, error) {
	key := "refresh_token:" + sessionID + ":" + email
	val, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
