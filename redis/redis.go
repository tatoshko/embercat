package redis

import (
	"gopkg.in/redis.v3"
	"log"
	"net/url"
	"os"
	"strings"
)

var REDIS_URL = os.Getenv("REDIS_URL")

func GetClient() *redis.Client {
	if REDIS_URL == "" {
		log.Printf("Redis can't be inited, because of empty HerokuURL")
		return nil
	}

	log.Printf("HerokuURL %s", REDIS_URL)

	var password string
	var addr string

	if !strings.Contains(REDIS_URL, "localhost") {
		parsedURL, _ := url.Parse(REDIS_URL)
		password, _ = parsedURL.User.Password()
		addr = parsedURL.Host
	}

	log.Printf("Connecting to ADDR: %s, PASWORD: %s", addr, password)

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}
