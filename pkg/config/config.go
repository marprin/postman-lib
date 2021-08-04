package config

import (
	"github.com/marprin/postman-lib/pkg/env"
	gcfg "gopkg.in/gcfg.v1"
)

func ReadModuleConfig(cfg interface{}, path string, module string) bool {
	environ := env.Get()
	fname := path + "/" + module + "." + environ + ".ini"
	err := gcfg.ReadFileInto(cfg, fname)
	if err == nil {
		return true
	}

	return false
}

func ReadModuleConfigWithErr(cfg interface{}, path string, module string) error {
	environ := env.Get()
	fname := path + "/" + module + "." + environ + ".ini"
	return gcfg.ReadFileInto(cfg, fname)
}
