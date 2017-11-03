package redis

import (
	"github.com/go-redis/redis"
)

// Engine module configuration.
type Engine struct {
	Client *redis.Client
}

// Get a specific value from a redis instance.
func (r *Engine) Get(key string) (interface{}, error) {
	return r.Client.Get(key).Result()
}

// Put specific value into a redis instance.
func (r *Engine) Put(key string, value interface{}) error {
	return r.Client.Set(key, value, 0).Err()
}
