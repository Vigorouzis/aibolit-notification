package http

type Config struct {
	Port string `env:"PORT"`
}

func (cfg Config) Addr() string {
	return ":" + cfg.Port
}
