package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

//func init() {
//	rdb = redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//}

func Set(key, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}

func Get(key string) (string, error) {
	return rdb.Get(ctx, "key").Result()
}
