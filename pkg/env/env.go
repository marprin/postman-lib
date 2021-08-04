package env

import (
	"os"
)

// Environment List
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

var env string

func init() {
	env = os.Getenv("STAGE")
	if env == "" {
		env = EnvDevelopment
	}
}

//Get return string of current environment flag
func Get() string {
	return env
}

// IsDevelopment check if current env is development
func IsDevelopment() bool {
	return EnvDevelopment == env
}

// IsProduction check if current env is production
func IsProduction() bool {
	return EnvProduction == env
}
