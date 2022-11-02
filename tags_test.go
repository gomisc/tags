package tags_test

import (
	"log"
	"os"
	"reflect"
	"time"

	"git.eth4.dev/golibs/errors"
	"git.eth4.dev/golibs/filepaths"

	"git.eth4.dev/golibs/tags"
)

type RedisConf struct {
	Endpoint string `env:"REDIS_ENDPOINT" default:"redis:6379"`
}

type Config struct {
	Endpoint                     string `env:"SERVICE_LISTEN_ENDPOINT" default:"0.0.0.0:8080"`
	TempDirectory                string `env:"SERVICE_TMP" default:"/tmp/servicename" checkfs:"mkdir"`
	Redis                        RedisConf
	StartTimeout                 time.Duration `env:"SERVICE_START_TIMEOUT" default:"15s"`
	LimitConnections             int           `env:"SERVICE_CONN_LIMIT" default:"1000"`
	RetryPolicyMaxAttempts       uint64        `env:"SERVICE_RETRY_MAX_ATTEMPTS" default:"10"`
	RetryPolicyBackoffMultiplier float64       `env:"SERVICE_RETRY_BACKOFF_MULTIPLIER" default:"1.5"`
}

type confSpec struct {
	Default tags.TagProcessor
	Env     tags.TagProcessor
	CheckFS tags.TagProcessor
}

func Example_tagSpecification() {
	spec, err := tags.ParseSpec(confSpec{
		Default: tags.DirectSetter(),
		Env:     tags.SetterFromEnv(),
		CheckFS: tags.TagProcFunc(func(f reflect.Value, tv string) error {
			if !filepaths.FileExists(f.String()) {
				switch tv {
				case "mkdir":
					return os.MkdirAll(f.String(), os.ModePerm)
				case "mkfile":
					fd, err := os.Create(f.String())
					if err != nil {
						return errors.Wrap(err, "create file")
					}

					return fd.Close()
				}
			}

			return nil
		}),
	})
	if err != nil {
		log.Fatalln(err)
	}

	conf := &Config{}
	if err = spec.Apply(conf); err != nil {
		log.Fatalln(err)
	}
}
