package config

const (
	ElasticSearch = "elasticsearch"
)

type Config struct {
	Resource Resource `toml:"resource"`
	Handler  Handler  `toml:"handler"`
}

type Resource struct {
	Pod     string `toml:"pod"`
	Service string `toml:"service"`
}

type Handler struct {
	Name   string        `toml:"handler"`
	Config HandlerConfig `toml:"handlerconfig"`
}

type HandlerConfig struct{}
