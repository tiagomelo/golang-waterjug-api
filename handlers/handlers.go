// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package handlers

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/tiagomelo/golang-waterjug-api/cache"
	v1 "github.com/tiagomelo/golang-waterjug-api/handlers/v1"
)

// ApiMuxConfig struct holds the configuration for the API.
type ApiMuxConfig struct {
	Cache cache.CacheService
	Log   *slog.Logger
}

// NewApiMux creates and returns a new mux.Router configured with version 1 (v1) routes.
func NewApiMux(c *ApiMuxConfig) *mux.Router {
	return v1.Routes(&v1.Config{
		Cache: c.Cache,
		Log:   c.Log,
	})
}
