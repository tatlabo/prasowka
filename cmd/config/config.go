package config

import "github.com/gofor-little/env"

type Config struct {
	Port string
	DNS  string
}

func (cf *Config) New() Config {
	if err := env.Load("././.env"); err != nil {
		panic(err)
	}
	cf.Port = env.Get("PORT", "8765")
	cf.DNS = env.Get("DNS", "./db/webservice.db")
	return *cf
}
