package conf

import "github.com/gofor-little/env"

type Config struct {
	Port string
	DNS  string
}

func (cf *Config) New() Config {
	if err := env.Load("././.env"); err != nil {
		panic(err)
	}
	cf.Port = env.Get("PORT", "8090")
	cf.DNS = env.Get("DNS", "./database/prasa.db")
	return *cf
}
