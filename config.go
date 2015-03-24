package imstor

import "github.com/deiwin/gonfigure"

var (
	rootPathEnvProperty = gonfigure.NewRequiredEnvProperty("IMSTOR_ROOT_PATH")
)

type Config struct {
	RootPath  string
	CopySizes []Size
}

func NewConfig(copySizes []Size) *Config {
	return &Config{
		RootPath:  rootPathEnvProperty.Value(),
		CopySizes: copySizes,
	}
}

// Size specifies a set of dimensions and a name that a copy of an image will
// be stored as
type Size struct {
	Name   string
	Height uint
	Width  uint
}
