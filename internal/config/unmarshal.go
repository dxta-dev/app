package config

import (
	"github.com/knadh/koanf/v2"
)

func unmarshal(k *koanf.Koanf) (Config, error) {
	out := config{}

	err := k.UnmarshalWithConf("", &out, koanf.UnmarshalConf{Tag: "koanf", FlatPaths: true})

	if err != nil {
		return Config{}, err
	}

	return getConfig(&out)
}
