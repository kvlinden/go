/*
Run this API integration test using the test utility option.
*/
package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
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
type WhereMap map[string]interface{}

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
		dateValue := *value
		filterConditions[fieldName] = dateValue.Format("2006-01-02")
	}
}
func BuildSearchConditions(searchString *string, searchConditions SearchConditions) {
	if searchString != nil {
		searchConditions["tableName.searchString"] = *searchString
		intValue, err := strconv.ParseInt(*searchString, 10, 64)
		if err == nil {
			searchConditions["tableName.SearchInt"] = intValue
		}
		floatValue, err := strconv.ParseFloat(*searchString, 64)
		if err == nil {
			searchConditions["tableName.SearchFloat"] = floatValue
		}
	}
}
func BuildFilterConditions(urlParameters *UrlParameters, filterConditions FilterConditions) {
	AddWhereIntCondition("tableName.IntFilter", urlParameters.IntFilter, filterConditions)
	AddWhereFloatCondition("tableName.FloatFilter", urlParameters.FloatFilter, filterConditions)
	AddWhereStringCondition("tableName.StringFilter", urlParameters.StringFilter, filterConditions)
	AddWhereBoolCondition("tableName.BoolFilter", urlParameters.BoolFilter, filterConditions)
	AddWhereDateCondition("tableName.DateFilter", urlParameters.DateFilter, filterConditions)
	AddWhereIntCondition("tableName.IntIgnoreFilter", urlParameters.IntIgnoreFilter, filterConditions)
	AddWhereFloatCondition("tableName.FloatIgnoreFilter", urlParameters.FloatIgnoreFilter, filterConditions)
	AddWhereStringCondition("tableName.StringIgnoreFilter", urlParameters.StringIgnoreFilter, filterConditions)
	AddWhereBoolCondition("tableName.BoolIgnoreFilter", urlParameters.BoolIgnoreFilter, filterConditions)
	AddWhereDateCondition("tableName.DateIgnoreFilter", urlParameters.DateIgnoreFilter, filterConditions)
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
}

// repository utilities
var FLOAT_EPSILON = 0.1

func Nameify(name string) string {
	return strings.Replace(name, ".", "_", -1)
}

func BuildWhereSearchConditionAndMap(whereConditionBuffer *bytes.Buffer, firstClause *bool, searchConditions *SearchConditions, whereMap WhereMap) {
	if len(*searchConditions) > 0 {
		for key, value := range *searchConditions {
			if *firstClause {
				whereConditionBuffer.WriteString("WHERE (")
				*firstClause = false
			} else {
				whereConditionBuffer.WriteString(" OR ")
			}
			switch value.(type) {
			case string:
				whereConditionBuffer.WriteString(fmt.Sprintf("%s LIKE '%%:%s%%'", key, Nameify(key)))
			case int64:
				whereConditionBuffer.WriteString(fmt.Sprintf("%s = :%s", key, Nameify(key)))
			case float64:
				whereConditionBuffer.WriteString(fmt.Sprintf("%s BETWEEN (:%s-%f) AND (:%s+%f)", key, Nameify(key), FLOAT_EPSILON, Nameify(key), FLOAT_EPSILON))
			}
			whereMap[Nameify(key)] = value
		}
		whereConditionBuffer.WriteString(") ")
	}
}

func BuildWhereFilterConditionAndMap(whereConditionBuffer *bytes.Buffer, firstClause *bool, filterConditions *FilterConditions, whereMap WhereMap) {
	for key, value := range *filterConditions {
		if *firstClause {
			whereConditionBuffer.WriteString("WHERE ")
			*firstClause = false
		} else {
			whereConditionBuffer.WriteString(" AND ")
		}
		var valueString string
		switch value.(type) {
		case string:
			valueString = fmt.Sprintf("':%s'", Nameify(key))
		case int64, float64, bool, time.Time:
			valueString = fmt.Sprintf(":%s", Nameify(key))
		}
		whereConditionBuffer.WriteString(fmt.Sprintf(" %s = %s ", key, valueString))
		whereMap[Nameify(key)] = value
	}
}

func BuildWhereConditionAndMap(searchConditions *SearchConditions, filterConditions *FilterConditions) (string, WhereMap) {
	var whereConditionBuffer bytes.Buffer
	whereMap := make(WhereMap)
	firstClause := true
	BuildWhereSearchConditionAndMap(&whereConditionBuffer, &firstClause, searchConditions, whereMap)
	BuildWhereFilterConditionAndMap(&whereConditionBuffer, &firstClause, filterConditions, whereMap)
	return whereConditionBuffer.String(), whereMap
}

func main() {
	// Controller code
	// Build the struct here; a map wouldn't support the API documentation (@apiUse) and would give away the DB structure.
	var urlParameters = UrlParameters{
		//SearchString: createString("str1"), // This string can only be compared with strings.
		SearchString: createString("1"), // This string can be compared with integers/floats.
		//SearchString:       nil, // The search string is either a (single) string or nil.
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
	searchConditions := make(SearchConditions)
	BuildSearchConditions(urlParameters.SearchString, searchConditions)
	filterConditions := make(FilterConditions)
	BuildFilterConditions(&urlParameters, filterConditions)
	fmt.Println(fmt.Sprintf("\nServices:\n%+v \n%+v", searchConditions, filterConditions))

	// Repository code
	whereCondition, whereMap := BuildWhereConditionAndMap(&searchConditions, &filterConditions)
	fmt.Println(fmt.Sprintf("\nRepositories:\n%s\n%+v", whereCondition, whereMap))

	testMap := map[string]interface{}{"tableName_FloatFilter": float64(1.0), "tableName_StringFilter": "str2", "tableName_SearchFloat": float64(1.0), "tableName_searchString": int64(1), "tableName_SearchInt": int64(1), "tableName_BoolFilter": true, "tableName_DateFilter": "2014-11-12", "tableName_IntFilter": int64(1)}
	fmt.Println(fmt.Sprintf("\nTesting\n%+v", testMap))

	for k, v := range whereMap {
		fmt.Println(fmt.Sprintf("m1*** %s: '%s' (%T)", k, v, v))
	}
	for k, v := range testMap {
		fmt.Println(fmt.Sprintf("m2*** %s: '%s' (%T)", k, v, v))
	}

}
