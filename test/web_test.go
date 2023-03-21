package test

import (
	"net/http"
	"testing"
)

func TestWeb(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("你好！"))
	})
	err := http.ListenAndServe(":8686", nil)
	if err != nil {
		return
	}
}
