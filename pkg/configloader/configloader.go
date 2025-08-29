package configloader

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"log"
	"strings"
)

type Option struct {
	Prefix       string
	Delimiter    string
	Separator    string
	YamlFilePath string
	CallbackEnv  func(string) string
}

func defaultCallbackEnv(source, prefix, separator string) string {
	base := strings.ToLower(strings.TrimPrefix(source, prefix))
	return strings.ReplaceAll(base, separator, ".")
}

func Load(options Option, config interface{}) error {
	k := koanf.New(options.Delimiter)

	if options.YamlFilePath != "" {
		if err := k.Load(file.Provider(options.YamlFilePath), yaml.Parser()); err != nil {
			log.Fatalf("Error loading config file: %v", err)
			return err
		}
	}

	callback := options.CallbackEnv
	if callback == nil {
		callback = func(source string) string {
			return defaultCallbackEnv(source, options.Prefix, options.Separator)
		}
	}

	if err := k.Load(env.Provider(options.Prefix, options.Delimiter, callback), nil); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
		return err
	}

	// Unmarshal into provided config structure (passing address)
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
		return err
	}

	return nil
}
