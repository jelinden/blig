package config

var Conf Config

func SetConfig(config Config) {
	Conf = config
}

type Config struct {
	AdminUsername *string
	AdminPassword *string
}
