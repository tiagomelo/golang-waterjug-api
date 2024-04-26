// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cache

import (
	"context"
	"time"
)

// CacheService defines the interface for a cache.
type CacheService interface {
	// Get retrieves the value associated with the given key from the cache.
	Get(ctx context.Context, key string) (string, error)
	// Set sets the value associated with the given key in the cache.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}
