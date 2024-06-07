package config

import (
	"context"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

var prefix = "DXTA_"

func Load(ctx context.Context, debug bool) (Config, error) {
	_ = k.Load(file.Provider("config.toml"), toml.Parser())

	if debug {
		_ = k.Load(file.Provider("config.dev.toml"), toml.Parser())
	}

	_ = k.Load(
		env.Provider(
			prefix,
			".",
			func(s string) string {
				return strings.Replace(
					strings.ToLower(strings.TrimPrefix(s, prefix)), "_", ".", -1)
			}),
		nil,
	)

	c, err := unmarshal(k)

	if (err != nil) {
		return Config{}, err
	}

	out, err := getConfig(ctx, &c)

	return out, nil
}
