// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package v1

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tiagomelo/golang-waterjug-api/cache"
	"github.com/tiagomelo/golang-waterjug-api/handlers/v1/waterjug"
	"github.com/tiagomelo/golang-waterjug-api/middleware"
)

// Config struct holds the database connection and logger.
type Config struct {
	Cache cache.CacheService
	Log   *slog.Logger
}

// Routes initializes and returns a new router with configured routes.
func Routes(c *Config) *mux.Router {
	router := mux.NewRouter()
	initializeRoutes(c.Cache, router)
	router.Use(
		func(h http.Handler) http.Handler {
			return middleware.Logger(c.Log, h)
		},
		middleware.Compress,
		middleware.PanicRecovery,
	)
	return router
}

// initializeRoutes sets up the routes.
func initializeRoutes(cache cache.CacheService, router *mux.Router) {
	waterjugHandlers := waterjug.New(cache)
	router.HandleFunc("/v1/measure", waterjugHandlers.Measure).Methods(http.MethodPost)
}
