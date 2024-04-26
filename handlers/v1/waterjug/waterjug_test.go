// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package waterjug

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/golang-waterjug-api/cache"
	"github.com/tiagomelo/golang-waterjug-api/measurement/models"
)

func TestMeasure(t *testing.T) {
	testCases := []struct {
		name                          string
		input                         string
		mockJsonDecode                func(r io.Reader, v any) error
		mockRetrieveSolutionFromCache func(ctx context.Context, cs cache.CacheService,
			newMeasurement *models.NewMeasurement) (*models.Solution, error)
		mockStoreSolutionInCache func(ctx context.Context, cs cache.CacheService,
			measurement *models.NewMeasurement, solution *models.Solution,
			expiration time.Duration) error
		expectedOutput     string
		expectedStatusCode int
	}{
		{
			name:  "happy path, no solution stored in cache yet",
			input: `{"x_capacity":2,"y_capacity":100,"z_amount_wanted":96}`,
			mockRetrieveSolutionFromCache: func(ctx context.Context, cs cache.CacheService, newMeasurement *models.NewMeasurement) (*models.Solution, error) {
				return nil, nil
			},
			mockStoreSolutionInCache: func(ctx context.Context, cs cache.CacheService, measurement *models.NewMeasurement, solution *models.Solution, expiration time.Duration) error {
				return nil
			},
			expectedOutput:     "{\"solution\":[{\"step\":1,\"bucketX\":0,\"bucketY\":100,\"action\":\"Fill bucket Y\"},{\"step\":2,\"bucketX\":2,\"bucketY\":98,\"action\":\"Transfer from bucket Y to X\"},{\"step\":3,\"bucketX\":0,\"bucketY\":98,\"action\":\"Empty bucket X\"},{\"step\":4,\"bucketX\":2,\"bucketY\":96,\"action\":\"Transfer from bucket Y to X\",\"status\":\"Solved\"}]}",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "happy path, solution stored in cache",
			input: `{"x_capacity":2,"y_capacity":100,"z_amount_wanted":96}`,
			mockRetrieveSolutionFromCache: func(ctx context.Context, cs cache.CacheService, newMeasurement *models.NewMeasurement) (*models.Solution, error) {
				return &models.Solution{
					Steps: []*models.Step{
						{
							Number:  1,
							BucketX: 1,
							BucketY: 1,
							Action:  "action",
							Status:  "status",
						},
					},
				}, nil
			},
			mockStoreSolutionInCache: func(ctx context.Context, cs cache.CacheService, measurement *models.NewMeasurement, solution *models.Solution, expiration time.Duration) error {
				return nil
			},
			expectedOutput:     "{\"solution\":[{\"step\":1,\"bucketX\":1,\"bucketY\":1,\"action\":\"action\",\"status\":\"status\"}]}",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "error when decoding payload",
			input: ``,
			mockJsonDecode: func(r io.Reader, v any) error {
				return errors.New("decode error")
			},
			expectedOutput:     "{\"error\":\"decode error\"}",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "input validation error",
			input:              `{}`,
			expectedOutput:     "{\"error\":\"[{\\\"field\\\":\\\"x_capacity\\\",\\\"error\\\":\\\"x_capacity is a required field\\\"},{\\\"field\\\":\\\"y_capacity\\\",\\\"error\\\":\\\"y_capacity is a required field\\\"},{\\\"field\\\":\\\"z_amount_wanted\\\",\\\"error\\\":\\\"z_amount_wanted is a required field\\\"}]\"}",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "no solution",
			input: `{"x_capacity":2,"y_capacity":6,"z_amount_wanted":5}`,
			mockRetrieveSolutionFromCache: func(ctx context.Context, cs cache.CacheService, newMeasurement *models.NewMeasurement) (*models.Solution, error) {
				return nil, nil
			},
			mockStoreSolutionInCache: func(ctx context.Context, cs cache.CacheService, measurement *models.NewMeasurement, solution *models.Solution, expiration time.Duration) error {
				return nil
			},
			expectedOutput:     "{\"error\":\"no solution\"}",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "error when retrieving solution from cache",
			input: `{"x_capacity":2,"y_capacity":100,"z_amount_wanted":96}`,
			mockRetrieveSolutionFromCache: func(ctx context.Context, cs cache.CacheService, newMeasurement *models.NewMeasurement) (*models.Solution, error) {
				return nil, errors.New("get error")
			},
			expectedOutput:     "{\"error\":\"get error\"}",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:  "error when storing solution in cache",
			input: `{"x_capacity":2,"y_capacity":100,"z_amount_wanted":96}`,
			mockRetrieveSolutionFromCache: func(ctx context.Context, cs cache.CacheService, newMeasurement *models.NewMeasurement) (*models.Solution, error) {
				return nil, nil
			},
			mockStoreSolutionInCache: func(ctx context.Context, cs cache.CacheService, measurement *models.NewMeasurement, solution *models.Solution, expiration time.Duration) error {
				return errors.New("set error")
			},
			expectedOutput:     "{\"error\":\"set error\"}",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	originalJsonDecode := jsonDecode
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				jsonDecode = originalJsonDecode
			}()
			if tc.mockJsonDecode != nil {
				jsonDecode = tc.mockJsonDecode
			}
			retrieveSolutionFromCache = tc.mockRetrieveSolutionFromCache
			storeSolutionInCache = tc.mockStoreSolutionInCache
			req, err := http.NewRequest(http.MethodPost, "measure", bytes.NewBuffer([]byte(tc.input)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			h := New(nil)
			handler := http.HandlerFunc((h).Measure)
			handler.ServeHTTP(recorder, req)
			require.Equal(t, tc.expectedStatusCode, recorder.Code)
			require.Equal(t, tc.expectedOutput, recorder.Body.String())
		})
	}
}
