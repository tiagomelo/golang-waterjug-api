//go:build integration

// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package integrationtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/golang-waterjug-api/cache"
	"github.com/tiagomelo/golang-waterjug-api/config"
	"github.com/tiagomelo/golang-waterjug-api/handlers"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	ctx := context.Background()
	cfg, err := config.ReadFromEnvFile(".env_test")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	redisCache, err := cache.NewRedisCache(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	apiMux := handlers.NewApiMux(
		&handlers.ApiMuxConfig{
			Cache: redisCache,
			Log:   log,
		},
	)
	testServer = httptest.NewServer(apiMux)
	defer testServer.Close()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestMeasurement(t *testing.T) {
	input := `{"x_capacity":2,"y_capacity":100,"z_amount_wanted":96}`
	expectedOutput := `{"solution":[{"step":1,"bucketX":0,"bucketY":100,"action":"Fill bucket Y"},{"step":2,"bucketX":2,"bucketY":98,"action":"Transfer from bucket Y to X"},{"step":3,"bucketX":0,"bucketY":98,"action":"Empty bucket X"},{"step":4,"bucketX":2,"bucketY":96,"action":"Transfer from bucket Y to X","status":"Solved"}]}`
	resp, err := http.Post(testServer.URL+"/v1/measure", "application/json", bytes.NewBuffer([]byte(input)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, string(b))
}
