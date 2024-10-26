package setting

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Setting .
type Setting struct {
	paths []string
}

// New Setting
func New(paths ...string) (s *Setting) {
	s = &Setting{
		paths: paths,
	}

	return
}

// Setup .
func (s *Setting) Setup(config interface{}) {
	for _, path := range s.paths {
		viper.SetConfigFile(path)

		if err := viper.MergeInConfig(); err != nil {
			fmt.Printf("Setup merge in config error %s", err.Error())
			panic(err)
		}
	}

	err := loadParseConfig(config)
	if err != nil {
		fmt.Printf("Load parse config error %s", err.Error())
		panic(err)
	}
}

func loadParseConfig(config interface{}) error {
	if err := mapstructure.Decode(viper.AllSettings(), config); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
