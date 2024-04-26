// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package doc

import "github.com/tiagomelo/golang-waterjug-api/measurement/models"

// swagger:route POST /v1/measure measure Get
// Get measurement.
// ---
// responses:
//		200: getMeasurementResponse
//		400: description: missing required fields
//		500: description: internal server error

// swagger:response getMeasurementResponse
type GetMeasurementResponseWrapper struct {
	// in:body
	Body models.NewMeasurement
}
