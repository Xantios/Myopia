package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnvHelper Make sure the env helper works
func TestEnvHelper(t *testing.T) {
	err := os.Setenv("TESTING", "SOMETHING")
	if err != nil {
		println(err.Error())
		t.FailNow()
	}

	if env("TESTING", "_") != "SOMETHING" {
		t.FailNow()
	}

	if env("ASSUMED_VALUE", "test") != "test" {
		t.FailNow()
	}
}

func Test_config_fails_hard_if_config_does_not_exist(t *testing.T) {
	assert.Panics(t, func() { GetConf("non_existent_file") }, nil)
}
