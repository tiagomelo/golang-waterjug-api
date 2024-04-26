// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/golang-waterjug-api/measurement/models"
)

func TestStoreSolution(t *testing.T) {
	testCases := []struct {
		name            string
		mockJsonMarshal func(v any) ([]byte, error)
		mockClosure     func(m *mockRedisCache)
		expectedError   error
	}{
		{
			name: "happy path",
			mockJsonMarshal: func(v any) ([]byte, error) {
				return []byte("some value"), nil
			},
			mockClosure: func(m *mockRedisCache) {},
		},
		{
			name: "error when serializing",
			mockJsonMarshal: func(v any) ([]byte, error) {
				return nil, errors.New("marshal error")
			},
			mockClosure:   func(m *mockRedisCache) {},
			expectedError: errors.New("failed to serialize solution: serializing solution: marshal error"),
		},
		{
			name: "error when setting value in Redis",
			mockJsonMarshal: func(v any) ([]byte, error) {
				return []byte("some value"), nil
			},
			mockClosure: func(m *mockRedisCache) {
				m.setErr = errors.New("set error")
			},
			expectedError: errors.New("set error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonMarshal = tc.mockJsonMarshal
			m := new(mockRedisCache)
			tc.mockClosure(m)
			err := StoreSolution(context.TODO(), m, &models.NewMeasurement{
				XCap:          2,
				YCap:          100,
				ZAmountWanted: 96,
			}, &models.Solution{}, 24*time.Hour)
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

func TestRetrieveSolution(t *testing.T) {
	testCases := []struct {
		name              string
		mockClosure       func(m *mockRedisCache)
		mockJsonUnmarshal func(data []byte, v any) error
		expectedError     error
	}{
		{
			name: "happy path",
			mockClosure: func(m *mockRedisCache) {
				m.val = "some cached val"
			},
			mockJsonUnmarshal: func(data []byte, v any) error {
				return nil
			},
		},
		{
			name: "error when getting value from Redis",
			mockClosure: func(m *mockRedisCache) {
				m.getErr = errors.New("get error")
			},
			expectedError: errors.New("failed to retrieve solution: get error"),
		},
		{
			name: "error when desserializing",
			mockClosure: func(m *mockRedisCache) {
				m.val = "some cached val"
			},
			mockJsonUnmarshal: func(data []byte, v any) error {
				return errors.New("unmarshal error")
			},
			expectedError: errors.New("deserializing solution: unmarshal error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonUnmarshal = tc.mockJsonUnmarshal
			m := new(mockRedisCache)
			tc.mockClosure(m)
			output, err := RetrieveSolution(context.TODO(), m, &models.NewMeasurement{})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error to occur, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error to occur, got nil`)
				}
				require.NotNil(t, output)
			}
		})
	}
}
