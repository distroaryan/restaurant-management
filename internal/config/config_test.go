package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()

	cfg := Load()


	defaultConfig := &Config{
		Env: "development",
		Port: 8080,
		MongoURI: "mongodb://127.0.0.1:27017",
	}

	assert.Equal(t, defaultConfig, cfg)
}