package extension

import (
	"github.com/redis/go-redis/v9"
	"github.com/suoaao/affordable-ai/pkg/conf"
)

var Redis *redis.Client

func init() {
	opt, err := redis.ParseURL(conf.RedisURL)
	if err != nil {
		panic(err)
	}
	Redis = redis.NewClient(opt)
}
