package main

import (
	"net/http"
	"testing"
)

func TestIndex(t *testing.T) {
	//exec.Command("./hello-plugin.exe", fmt.Sprintf("-port=%d", 30002))

	_, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

}
