package conf

import (
	"github.com/kelseyhightower/envconfig"
)

type conf struct {
	AdminAuthToken string   `required:"true" envconfig:"ADMIN_AUTH_TOKEN"`
	AuthTokens     []string `envconfig:"AUTH_TOKENS"`
	RedisURL       string   `envconfig:"REDIS_URL"`
	GhuToken       string   `envconfig:"GHU_TOKEN"`
}

func (c *conf) Auth(authToken string) bool {
	if authToken == c.AdminAuthToken {
		return true
	}
	for _, t := range c.AuthTokens {
		if authToken == t {
			return true
		}
	}
	return false
}

var Conf conf

func init() {
	err := envconfig.Process("", &Conf)
	if err != nil {
		panic(err)
	}
}
