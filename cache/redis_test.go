// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestNewRedisCache(t *testing.T) {
	testCases := []struct {
		name               string
		mockRedisNewClient func(opt *redis.Options) *redis.Client
		mockPing           func(ctx context.Context, client *redis.Client) error
		expectedError      error
	}{
		{
			name: "happy path",
			mockRedisNewClient: func(opt *redis.Options) *redis.Client {
				return new(redis.Client)
			},
			mockPing: func(ctx context.Context, client *redis.Client) error {
				return nil
			},
		},
		{
			name: "ping error",
			mockRedisNewClient: func(opt *redis.Options) *redis.Client {
				return new(redis.Client)
			},
			mockPing: func(ctx context.Context, client *redis.Client) error {
				return errors.New("ping error")
			},
			expectedError: errors.New("pinging redis instance: ping error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redisNewClient = tc.mockRedisNewClient
			ping = tc.mockPing
			r, err := NewRedisCache(context.TODO(), "host", "port")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error to occur, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error to occur, got nil`)
				}
				require.NotNil(t, r)
			}
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name           string
		mockClosure    func(mSCmd *mockRedisStringCmd)
		expectedOutput string
		expectedError  error
	}{
		{
			name: "happy path",
			mockClosure: func(mSCmd *mockRedisStringCmd) {
				mSCmd.val = "some cached value"
			},
			expectedOutput: "some cached value",
		},
		{
			name: "key does not exists",
			mockClosure: func(mSCmd *mockRedisStringCmd) {
				mSCmd.err = redis.Nil
			},
		},
		{
			name: "error",
			mockClosure: func(mSCmd *mockRedisStringCmd) {
				mSCmd.err = errors.New("get error")
			},
			expectedError: errors.New(`getting key "some key": get error`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mRedisStringCmd := new(mockRedisStringCmd)
			tc.mockClosure(mRedisStringCmd)
			mRedisClient := new(mockRedisClient)
			mRedisClient.redisStringCmd = mRedisStringCmd
			rc := &redisCache{
				redisClient: mRedisClient,
			}
			output, err := rc.Get(context.TODO(), "some key")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error to occur, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error to occur, got nil`)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestSet(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(mStCmd *mockRedisStatusCmd)
		expectedError error
	}{
		{
			name:        "happy path",
			mockClosure: func(mStCmd *mockRedisStatusCmd) {},
		},
		{
			name: "error",
			mockClosure: func(mStCmd *mockRedisStatusCmd) {
				mStCmd.err = errors.New("set error")
			},
			expectedError: errors.New(`set key "some key" and expiration 24h0m0s: set error`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mRedisStCmd := new(mockRedisStatusCmd)
			tc.mockClosure(mRedisStCmd)
			mRedisClient := new(mockRedisClient)
			mRedisClient.redisStatusCmd = mRedisStCmd
			rc := &redisCache{
				redisClient: mRedisClient,
			}
			err := rc.Set(context.TODO(), "some key", "some value", 24*time.Hour)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error to occur, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error to occur, got nil`)
				}
			}
		})
	}
}
