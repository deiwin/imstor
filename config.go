package imstor

import "github.com/deiwin/gonfigure"

var (
	rootPathEnvProperty = gonfigure.NewRequiredEnvProperty("IMSTOR_ROOT_PATH")
)

type Config struct {
	RootPath string
}

func NewConfig() *Config {
	return &Config{
		RootPath: rootPathEnvProperty.Value(),
	}
}
