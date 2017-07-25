/*
Run this API integration test using the test utility option.
*/
package main

import (
	"fmt"
	"reflect"
	"time"
)

// Controller utilities
// Get/verify URL parameter values and build parameters struct.

type UrlSearchParameters struct {
	SearchString string
	SearchInt    int64
	SearchFloat  float64
}
type UrlParameters struct {
	SearchString       *UrlSearchParameters
	IntFilter          *int64
	FloatFilter        *float64
	StringFilter       *string
	BoolFilter         *bool
	DateFilter         *time.Time
	IntIgnoreFilter    *int64
	FloatIgnoreFilter  *float64
	StringIgnoreFilter *string
	BoolIgnoreFilter   *bool
	DateIgnoreFilter   *time.Time
}

func createInt64(x int64) *int64 {
	return &x
}
func createFloat64(x float64) *float64 {
	return &x
}
func createString(x string) *string {
	return &x
}
func createBool(x bool) *bool {
	return &x
}
func createDate(dateString string) *time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", dateString)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	return &t
}

// Service utilities
// Convert URL parameter structure to where condition map, ignoring nil values.

type WhereConditions map[string]interface{}

// repository utilities

func main() {
	// Controller code
	var urlParameters = UrlParameters{
		SearchString: &UrlSearchParameters{
			SearchString: "str1",
			SearchInt:    1,
			SearchFloat:  1.0,
		},
		IntFilter:          createInt64(1),
		FloatFilter:        createFloat64(1.0),
		StringFilter:       createString("str2"),
		BoolFilter:         createBool(true),
		DateFilter:         createDate("2014-11-12T11:45:26.371Z"),
		IntIgnoreFilter:    nil,
		FloatIgnoreFilter:  nil,
		StringIgnoreFilter: nil,
		BoolIgnoreFilter:   nil,
		DateIgnoreFilter:   nil,
	}
	fmt.Println(fmt.Sprintf("Controller: creates %+v", urlParameters))

	// Services code
	whereConditions := make(WhereConditions)
	v := reflect.ValueOf(urlParameters)
	for i := 0; i < v.NumField(); i++ {
		fmt.Println(v.Type().Field(i).Name, v.Field(i).Type().String(), v.Field(i).Interface())
		fmt.Printf("(%v, %T)\n", v.Field(i).Interface(), v.Field(i).Interface())
		if v.Field(i).Interface() != nil {
			switch ptr := v.Field(i).Interface().(type) {
			case *string:
				fmt.Println(*ptr)
			}
		}
	}
	fmt.Println(whereConditions)

	// Repository code
}
