// Package config manages the configuration to be 12factor.
package config

import "os"

// LookupEnvOrString returns the environmnental variable behind the key
// if the key is not found, it returns defaultVal.
func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
