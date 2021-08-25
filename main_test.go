package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnvHelper Make sure the env helper works
func TestEnvHelper(t *testing.T) {
	t.Setenv("TESTING", "SOMETHING")

	if env("TESTING", "_") != "SOMETHING" {
		t.FailNow()
	}

	if env("ASSUMED_VALUE", "test") != "test" {
		t.FailNow()
	}
}

func TestConfigFailsHardIfConfigDoesNotExist(t *testing.T) {
	assert.Panics(t, func() { GetConf("non_existent_file") }, nil)
}

func TestConfigFileExists(t *testing.T) {
	GetConf("config.test.yaml")
}
