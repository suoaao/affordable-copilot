package conf

import (
	"github.com/kelseyhightower/envconfig"
)

type copilotConf struct {
	GhuToken string `envconfig:"GHU_TOKEN"`
}

type openaiConf struct {
	ApiKey string `envconfig:"OPENAI_KEY"`
}

type conf struct {
	AuthToken string `required:"true" envconfig:"AUTH_TOKEN"`
	RedisURL  string `envconfig:"REDIS_URL"`
	Copilot   copilotConf
	Openai    openaiConf
}

var Conf conf

func init() {
	err := envconfig.Process("", &Conf)
	if err != nil {
		panic(err)
	}
}
