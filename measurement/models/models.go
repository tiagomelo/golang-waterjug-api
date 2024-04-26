// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package models

// NewMeasurement represents the measurements for the water jug problem.
type NewMeasurement struct {
	XCap          int `json:"x_capacity" validate:"required,gt=0,lt=10000"`      // XCap represents the capacity of jug X.
	YCap          int `json:"y_capacity" validate:"required,gt=0,lt=10000"`      // YCap represents the capacity of jug Y.
	ZAmountWanted int `json:"z_amount_wanted" validate:"required,gt=0,lt=10000"` // ZAmountWanted represents the desired amount of water Z.
}

// Step represents a step in the solution to the water jug problem.
type Step struct {
	Number  int    `json:"step"`             // Number represents the step number.
	BucketX int    `json:"bucketX"`          // BucketX represents the amount of water in jug X.
	BucketY int    `json:"bucketY"`          // BucketY represents the amount of water in jug Y.
	Action  string `json:"action"`           // Action represents the action taken in this step.
	Status  string `json:"status,omitempty"` // Status represents the status of the step, if applicable.
}

// Solution represents the solution to the water jug problem.
type Solution struct {
	Steps []*Step `json:"solution"` // Steps is a slice of steps representing the solution path.
}
