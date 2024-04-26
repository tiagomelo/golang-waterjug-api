// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package measurement

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/golang-waterjug-api/measurement/models"
)

func TestMeasure(t *testing.T) {
	testCases := []struct {
		name           string
		xMax           int
		yMax           int
		target         int
		expectedOutput *models.Solution
	}{
		{
			name:   "measurement is possible",
			xMax:   2,
			yMax:   100,
			target: 96,
			expectedOutput: &models.Solution{
				Steps: []*models.Step{
					{
						Number:  1,
						BucketX: 0,
						BucketY: 100,
						Action:  "Fill bucket Y",
					},
					{
						Number:  2,
						BucketX: 2,
						BucketY: 98,
						Action:  "Transfer from bucket Y to X",
					},
					{
						Number:  3,
						BucketX: 0,
						BucketY: 98,
						Action:  "Empty bucket X",
					},
					{
						Number:  4,
						BucketX: 2,
						BucketY: 96,
						Action:  "Transfer from bucket Y to X",
						Status:  "Solved",
					},
				},
			},
		},
		{
			name:   "different gcd",
			xMax:   8,
			yMax:   12,
			target: 5,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := Measure(tc.xMax, tc.yMax, tc.target)
			require.Equal(t, tc.expectedOutput, output)
		})
	}
}
