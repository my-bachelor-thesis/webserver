package redis

import (
	"context"
	"encoding/json"
	"errors"
	goredis "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	"webserver/internal/config"
)

var (
	ctx = context.Background()
	rh  *rejson.Handler
	cli *goredis.Client
)

func Connect() error {
	rh = rejson.NewReJSONHandler()
	cli = goredis.NewClient(&goredis.Options{Addr: config.GetInstance().RedisURL})
	rh.SetGoRedisClient(cli)
	var testValue struct {
		Value string
	}
	testValue.Value = "test"
	return Set("test_key", testValue)
}

func CloseConnection() {
	_ = cli.FlushAll(ctx).Err()
	_ = cli.Close()
}

// Set - value should be a pointer to a struct
func Set(key string, value interface{}) error {
	res, err := rh.JSONSet(key, ".", value)
	if err != nil {
		return err
	}

	if res.(string) != "OK" {
		return errors.New("failed to set")
	}
	return nil
}

// Get - structToFill has to be a pointer to a struct
func Get(key string, structToFill interface{}) error {
	bytes, err := redis.Bytes(rh.JSONGet(key, "."))
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, structToFill)
}

func IsErrKeyNotFound(err error) bool {
	return err.Error() == "redis: nil"
}
