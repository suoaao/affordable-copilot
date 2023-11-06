package extension

import (
	"github.com/redis/go-redis/v9"
	"github.com/suoaao/affordable-copilot/pkg/conf"
)

var Redis *redis.Client

func init() {
	opt, err := redis.ParseURL(conf.Conf.RedisURL)
	if err != nil {
		panic(err)
	}
	Redis = redis.NewClient(opt)
}
