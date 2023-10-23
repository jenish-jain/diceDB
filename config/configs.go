package config

type Configs struct {
	Host string `default:"0.0.0.0"`
	Port int    `default:"7379"`
}

func InitConfigWithDefaultValues() *Configs {
	return &Configs{}
}
