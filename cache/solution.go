// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/tiagomelo/golang-waterjug-api/measurement/models"
)

// For ease of unit testing.
var (
	jsonMarshal   = json.Marshal
	jsonUnmarshal = json.Unmarshal
)

// serializeSolution converts a Solution instance into a JSON string.
func serializeSolution(solution *models.Solution) (string, error) {
	jsonBytes, err := jsonMarshal(solution)
	if err != nil {
		return "", errors.Wrap(err, "serializing solution")
	}
	return string(jsonBytes), nil
}

// deserializeSolution converts a JSON string back into a Solution instance.
func deserializeSolution(data string) (*models.Solution, error) {
	var solution models.Solution
	err := jsonUnmarshal([]byte(data), &solution)
	if err != nil {
		return nil, errors.Wrap(err, "deserializing solution")
	}
	return &solution, nil
}

// StoreSolution stores the serialized Solution in the cache.
func StoreSolution(ctx context.Context, cache CacheService, measurement *models.NewMeasurement, solution *models.Solution, expiration time.Duration) error {
	serializedSolution, err := serializeSolution(solution)
	if err != nil {
		return errors.Wrap(err, "failed to serialize solution")
	}
	return cache.Set(ctx, solutionCacheKey(measurement), serializedSolution, expiration)
}

// RetrieveSolution retrieves the serialized Solution from the cache and deserializes it.
func RetrieveSolution(ctx context.Context, cache CacheService, measurement *models.NewMeasurement) (*models.Solution, error) {
	serializedSolution, err := cache.Get(ctx, solutionCacheKey(measurement))
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve solution")
	}
	if serializedSolution == "" {
		return nil, nil
	}
	return deserializeSolution(serializedSolution)
}

// solutionCacheKey generates a cache key based on the measurement parameters.
func solutionCacheKey(measurement *models.NewMeasurement) string {
	return fmt.Sprintf("%d#%d#%d", measurement.XCap, measurement.YCap, measurement.ZAmountWanted)
}
