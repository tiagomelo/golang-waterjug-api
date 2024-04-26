// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cache

import (
	"context"
	"time"
)

type mockRedisStringCmd struct {
	err error
	val string
}

func (m *mockRedisStringCmd) Err() error {
	return m.err
}

func (m *mockRedisStringCmd) Val() string {
	return m.val
}

type mockRedisStatusCmd struct {
	err error
}

func (m *mockRedisStatusCmd) Err() error {
	return m.err
}

type mockRedisClient struct {
	redisStringCmd *mockRedisStringCmd
	redisStatusCmd *mockRedisStatusCmd
}

func (m *mockRedisClient) Get(ctx context.Context, key string) redisStringCmd {
	return m.redisStringCmd
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) redisStatusCmd {
	return m.redisStatusCmd
}

type mockRedisCache struct {
	val    string
	getErr error
	setErr error
}

func (m *mockRedisCache) Get(ctx context.Context, key string) (string, error) {
	return m.val, m.getErr
}

func (m *mockRedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.setErr
}
