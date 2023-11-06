package conf

import (
	"github.com/kelseyhightower/envconfig"
)

type conf struct {
	AuthToken string `required:"true" envconfig:"AUTH_TOKEN"`
	RedisURL  string `envconfig:"REDIS_URL"`
	GhuToken  string `envconfig:"GHU_TOKEN"`
}

var Conf conf

func init() {
	err := envconfig.Process("", &Conf)
	if err != nil {
		panic(err)
	}
}
