package conf

import "os"

var (
	AuthToken = os.Getenv("AUTH_TOKEN")
	GhuToken  = os.Getenv("GHU_TOKEN")
	RedisURL  = os.Getenv("REDIS_URL")
	OpenaiKey = os.Getenv("OPENAI_KEY")
)

func init() {
	if len(AuthToken) < 20 {
		panic("AUTH_TOKEN is invalid")
	}
	if len(GhuToken) < 20 {
		panic("GhuToken is invalid")
	}
	if len(RedisURL) < 10 {
		panic("RedisURL is invalid")
	}
	if len(OpenaiKey) < 10 {
		panic("OpenaiKey is invalid")
	}
}
