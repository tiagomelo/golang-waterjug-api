// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// Config holds all configuration needed by this app.
type Config struct {
	RedisHost string `envconfig:"REDIS_HOST" required:"true"`
	RedisPort string `envconfig:"REDIS_PORT" required:"true"`
}

// For ease of unit testing.
var (
	godotenvLoad     = godotenv.Load
	envconfigProcess = envconfig.Process
)

// Read reads configuration from environment variables.
// It assumes that an '.env' file is present at current path.
func Read() (*Config, error) {
	if err := godotenvLoad(); err != nil {
		return nil, errors.Wrap(err, "loading env vars")
	}
	config := new(Config)
	if err := envconfigProcess("", config); err != nil {
		return nil, errors.Wrap(err, "processing env vars")
	}
	return config, nil
}

// ReadFromEnvFile reads configuration from the specified environment file.
func ReadFromEnvFile(envFilePath string) (*Config, error) {
	if err := godotenvLoad(envFilePath); err != nil {
		return nil, errors.Wrapf(err, "loading env vars from %s", envFilePath)
	}
	config := new(Config)
	if err := envconfigProcess("", config); err != nil {
		return nil, errors.Wrap(err, "processing env vars")
	}
	return config, nil
}
