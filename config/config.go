package config

import (
	"os"

	"github.com/orkungursel/go-eco"
)

func New() *Config {
	c := &Config{}
	if err := eco.Unmarshal(c); err != nil {
		panic(err)
	}

	return c
}

// GetProfile returns current mode by reading environment variable.
func (c *Config) GetProfile() string {
	profile, ok := os.LookupEnv("ACTIVE_PROFILE")
	if !ok {
		profile = "local"
	}

	return profile
}

// IsLocal returns true if current profile is local.
func (c *Config) IsLocal() bool {
	return c.GetProfile() == "local"
}

// IsProduction returns true if current profile is production.
func (c *Config) IsProduction() bool {
	return c.GetProfile() == "production"
}
