// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package validate

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/tiagomelo/golang-waterjug-api/measurement/models"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

const (
	// splitCount determines the number of substrings to return.
	splitCount = 2

	lessThanXAndYCapacitiesStructTagName = "less_than_x_y_cap"
)

// registerTranslationForLessThanXAndYCapacitiesStructTagName registers custom translation message
// when "less_than_x_y_cap" validation is violated.
func registerTranslationForLessThanXAndYCapacitiesStructTagName(ut ut.Translator) error {
	return ut.Add(lessThanXAndYCapacitiesStructTagName, "{0} 'z_amount_wanted' cannot be greater than x_capacity AND y_capacity", true)
}

// translationForLessThanXAndYCapacitiesStructTagName formats the message to be displayed
// for "less_than_x_y_cap" struct tag validation.
func translationForLessThanXAndYCapacitiesStructTagName(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(lessThanXAndYCapacitiesStructTagName, fe.Field(), fmt.Sprintf("%v", fe.Value()))
	return t
}

func init() {
	// Instantiate a validator.
	validate = validator.New()

	// Create a translator for english so the error messages are
	// more human-readable than technical.
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the english error messages for use.
	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		fmt.Println("error registering default transations:", err)
		os.Exit(1)
	}

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", splitCount)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// registers custom translation message when "less_than_x_y_cap" error tag is reported
	if err := validate.RegisterTranslation(lessThanXAndYCapacitiesStructTagName, translator, registerTranslationForLessThanXAndYCapacitiesStructTagName, translationForLessThanXAndYCapacitiesStructTagName); err != nil {
		fmt.Printf("error registering translations for %s tag: %v", lessThanXAndYCapacitiesStructTagName, err)
		os.Exit(1)
	}

	// registers validation for person.Person struct
	validate.RegisterStructValidation(NewMeasurementStructLevelValidation, models.NewMeasurement{})
}

// Check validates the provided model against it's declared tags.
func Check(val any) error {
	if err := validate.Struct(val); err != nil {
		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(translator),
			}
			fields = append(fields, field)
		}
		return fields
	}
	return nil
}

// NewMeasurementStructLevelValidation registers validation for models.NewMeasurement struct.
func NewMeasurementStructLevelValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(models.NewMeasurement)
	if req.ZAmountWanted > req.XCap && req.ZAmountWanted > req.YCap {
		sl.ReportError(nil, "", "", lessThanXAndYCapacitiesStructTagName, "")
	}
}
