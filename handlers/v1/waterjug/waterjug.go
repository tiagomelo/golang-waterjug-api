// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package waterjug

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/tiagomelo/golang-waterjug-api/cache"
	"github.com/tiagomelo/golang-waterjug-api/measurement"
	"github.com/tiagomelo/golang-waterjug-api/measurement/models"
	"github.com/tiagomelo/golang-waterjug-api/validate"
	"github.com/tiagomelo/golang-waterjug-api/web"
)

// handlers represents HTTP handlers for the water jug measurement service.
type handlers struct {
	cache cache.CacheService
}

// cache expiration time for cached solutions (24 hours).
const CACHE_EXPIRATION_24H = 24 * time.Hour

// For ease of unit testing.
var (
	// jsonDecode decodes a JSON request body into a given struct.
	jsonDecode = func(r io.Reader, v any) error {
		return json.NewDecoder(r).Decode(v)
	}
	// retrieveSolutionFromCache retrieves a solution from the cache based on the given measurement.
	retrieveSolutionFromCache = func(ctx context.Context,
		cs cache.CacheService, newMeasurement *models.NewMeasurement) (*models.Solution, error) {
		return cache.RetrieveSolution(ctx, cs, newMeasurement)
	}
	// storeSolutionInCache stores a solution in the cache.
	storeSolutionInCache = func(ctx context.Context,
		cs cache.CacheService, measurement *models.NewMeasurement,
		solution *models.Solution, expiration time.Duration) error {
		return cache.StoreSolution(ctx, cs, measurement, solution, expiration)
	}
)

// New creates a new handlers instance with the provided cache service.
func New(cache cache.CacheService) *handlers {
	return &handlers{
		cache: cache,
	}
}

// Measure is an HTTP handler for measuring water jug solutions.
func (h *handlers) Measure(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var newMeasurement models.NewMeasurement
	if err := jsonDecode(r.Body, &newMeasurement); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Check(newMeasurement); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	cachedSolution, err := retrieveSolutionFromCache(r.Context(), h.cache, &newMeasurement)
	if err != nil {
		web.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if cachedSolution != nil {
		web.RespondWithJson(w, http.StatusOK, cachedSolution)
		return
	}
	solution := measurement.Measure(newMeasurement.XCap, newMeasurement.YCap, newMeasurement.ZAmountWanted)
	if solution == nil {
		web.RespondWithError(w, http.StatusBadRequest, errors.New("no solution").Error())
		return
	}
	if err := storeSolutionInCache(r.Context(), h.cache, &newMeasurement, solution, CACHE_EXPIRATION_24H); err != nil {
		web.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	web.RespondWithJson(w, http.StatusOK, solution)
}
