// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// For ease of unit testing.
var (
	redisNewClient = redis.NewClient
	ping           = func(ctx context.Context, client *redis.Client) error {
		return client.Ping(ctx).Err()
	}
)

// redisStringCmd is an interface for a Redis string command.
type redisStringCmd interface {
	Err() error
	Val() string
}

// redisStringCmdWrapper is a wrapper for Redis string commands.
type redisStringCmdWrapper struct {
	sCmd *redis.StringCmd
}

// Err returns the error status of the wrapped command.
func (w *redisStringCmdWrapper) Err() error {
	return w.sCmd.Err()
}

// Val returns the string value of the wrapped command.
func (w *redisStringCmdWrapper) Val() string {
	return w.sCmd.Val()
}

// redisStatusCmd is an interface for a Redis status command.
type redisStatusCmd interface {
	Err() error
}

// redisStatusCmdWrapper is a wrapper for Redis status commands.
type redisStatusCmdWrapper struct {
	stCmd *redis.StatusCmd
}

// Err returns the error status of the wrapped command.
func (w *redisStatusCmdWrapper) Err() error {
	return w.stCmd.Err()
}

// redisClient is an interface for a Redis client.
type redisClient interface {
	// Get gets the value for the specified key.
	Get(ctx context.Context, key string) redisStringCmd
	// Set sets the value and expiration of a key.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) redisStatusCmd
}

// redisClientWrapper is a wrapper for the Redis client.
type redisClientWrapper struct {
	client *redis.Client
}

func (c *redisClientWrapper) Get(ctx context.Context, key string) redisStringCmd {
	return &redisStringCmdWrapper{c.client.Get(ctx, key)}
}

func (c *redisClientWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) redisStatusCmd {
	return &redisStatusCmdWrapper{c.client.Set(ctx, key, value, expiration)}
}

// redisCache represents a Redis cache.
type redisCache struct {
	redisClient redisClient
}

func (rc *redisCache) Get(ctx context.Context, key string) (string, error) {
	getCmd := rc.redisClient.Get(ctx, key)
	if err := getCmd.Err(); err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", errors.Wrapf(err, `getting key "%s"`, key)
	}
	return getCmd.Val(), nil
}

func (rc *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	statusCmd := rc.redisClient.Set(ctx, key, value, expiration)
	if err := statusCmd.Err(); err != nil {
		return errors.Wrapf(err, `set key "%v" and expiration %v`, key, expiration)
	}
	return nil
}

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(ctx context.Context, host, port string) (CacheService, error) {
	client := redisNewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	wrappedClient := &redisClientWrapper{client}
	if err := ping(ctx, client); err != nil {
		return nil, errors.Wrap(err, "pinging redis instance")
	}
	return &redisCache{wrappedClient}, nil
}
