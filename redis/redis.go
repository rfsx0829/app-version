package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

// Client redis.Client
var Client *redis.Client

// InitClient init client
func InitClient(host string, port int, password string, DB int) {
	opt := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       DB,
	}

	Client = redis.NewClient(&opt)
}

// GetHDataJSON returns json data of a hash table
func GetHDataJSON(name string) ([]byte, error) {
	mp, err := Client.HGetAll(name).Result()
	if err != nil {
		return nil, err
	}

	return json.Marshal(mp)
}

// GetSDataJSON returns json data of a set
func GetSDataJSON(name string) ([]byte, error) {
	sets, err := Client.SMembers(name).Result()
	if err != nil {
		return nil, err
	}

	return json.Marshal(sets)
}
