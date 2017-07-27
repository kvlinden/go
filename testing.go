package main

import "fmt"

func main() {

	m1 := map[string]interface{}{"city": "Johannesburg"}
	fmt.Println(fmt.Sprintf("%+v (%T)", m1, m1))
	fmt.Println(fmt.Sprintf("'%s' (%T)", m1["city"], m1["city"]))
	for k, v := range m1 {
		fmt.Println(fmt.Sprintf("%s: '%s' (%T)", k, v, v))
	}

	m2 := make(map[string]interface{})
	m2["city"] = "Johannesburg"
	fmt.Println(fmt.Sprintf("%+v (%T)", m2, m2))
	fmt.Println(fmt.Sprintf("'%s' (%T)", m2["city"], m2["city"]))
	for k, v := range m2 {
		fmt.Println(fmt.Sprintf("%s: '%s' (%T)", k, v, v))
	}
}
