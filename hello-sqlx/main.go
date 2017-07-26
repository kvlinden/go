/*
Run this API integration test using the test utility option.
*/
package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Controller utilities
// Get/verify URL parameter values and build parameters struct.

type UrlParameters struct {
	SearchString       *string
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

type SearchConditions map[string]interface{}
type FilterConditions map[string]interface{}

// Is there a way to handle these using *interface{}? I couldn't get it to work.
func AddWhereIntCondition(fieldName string, value *int64, filterConditions FilterConditions) {
	if value != nil {
		filterConditions[fieldName] = *value
	}
}
func AddWhereFloatCondition(fieldName string, value *float64, filterConditions FilterConditions) {
	if value != nil {
		filterConditions[fieldName] = *value
	}
}
func AddWhereStringCondition(fieldName string, value *string, filterConditions FilterConditions) {
	if value != nil {
		filterConditions[fieldName] = *value
	}
}
func AddWhereBoolCondition(fieldName string, value *bool, filterConditions FilterConditions) {
	if value != nil {
		filterConditions[fieldName] = *value
	}
}
func AddWhereDateCondition(fieldName string, value *time.Time, filterConditions FilterConditions) {
	if value != nil {
		filterConditions[fieldName] = *value
	}
}

// repository utilities

func main() {
	// Controller code
	var urlParameters = UrlParameters{
		SearchString: createString("str1"),
		//SearchString:       nil,  // The search string is either a (single) string or nil.
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
	fmt.Println(fmt.Sprintf("\nController:\n%+v", urlParameters))

	// Services code
	filterConditions := make(FilterConditions)
	searchConditions := make(SearchConditions)
	if urlParameters.SearchString != nil {
		searchConditions["SearchString"] = *urlParameters.SearchString
		intValue, err := strconv.ParseInt(*urlParameters.SearchString, 10, 64)
		if err == nil {
			searchConditions["SearchInt"] = intValue
		}
		floatValue, err := strconv.ParseFloat(*urlParameters.SearchString, 64)
		if err == nil {
			searchConditions["SearchFloat"] = floatValue
		}
	}
	AddWhereIntCondition("IntFilter", urlParameters.IntFilter, filterConditions)
	AddWhereFloatCondition("FloatFilter", urlParameters.FloatFilter, filterConditions)
	AddWhereStringCondition("StringFilter", urlParameters.StringFilter, filterConditions)
	AddWhereBoolCondition("BoolFilter", urlParameters.BoolFilter, filterConditions)
	AddWhereDateCondition("DateFilter", urlParameters.DateFilter, filterConditions)
	AddWhereIntCondition("IntIgnoreFilter", urlParameters.IntIgnoreFilter, filterConditions)
	AddWhereFloatCondition("FloatIgnoreFilter", urlParameters.FloatIgnoreFilter, filterConditions)
	AddWhereStringCondition("StringIgnoreFilter", urlParameters.StringIgnoreFilter, filterConditions)
	AddWhereBoolCondition("BoolIgnoreFilter", urlParameters.BoolIgnoreFilter, filterConditions)
	AddWhereDateCondition("DateIgnoreFilter", urlParameters.DateIgnoreFilter, filterConditions)
	// Generalizing this is probably possible using Reflection, but complicated. This code doesn't work properly.
	//v := reflect.ValueOf(urlParameters)
	//for i := 0; i < v.NumField(); i++ {
	//	fmt.Println(v.Type().Field(i).Name, v.Field(i).Type().String(), v.Field(i).Interface())
	//	fmt.Printf("(%v, %T)\n", v.Field(i).Interface(), v.Field(i).Interface())
	//	if v.Field(i).Interface() != nil {
	//		switch ptr := v.Field(i).Interface().(type) {
	//		case *string:
	//			fmt.Println(*ptr)
	//		}
	//	}
	//}
	fmt.Println(fmt.Sprintf("\nServices:\n%+v \n%+v", searchConditions, filterConditions))

	// Repository code
	var whereConditionBuffer bytes.Buffer
	//for key, value := range searchConditions {
	//	whereConditionBuffer.WriteString(fmt.Sprintf("OR %s = %s ", key, valueString))
	//}
	for key, value := range filterConditions {
		var valueString string
		switch v := value.(type) {
		case string:
			valueString = fmt.Sprintf("'%s'", v)
		case int64:
			valueString = fmt.Sprintf("%d", v)
		case float64:
			valueString = fmt.Sprintf("%f", v)
		case bool:
			valueString = strconv.FormatBool(v)
		case time.Time:
			valueString = fmt.Sprintf("'%s'", v.String())
		}
		whereConditionBuffer.WriteString(fmt.Sprintf("AND %s = %s ", key, valueString))
	}
	fmt.Println(fmt.Sprintf("\nRepositories:\n%s", whereConditionBuffer.String()))
}
