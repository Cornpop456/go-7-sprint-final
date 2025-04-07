package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		// пока сравнивать не будем, а просто выведем ответы
		// удалите потом этот вывод
		fmt.Println(response.Body.String())
	}
}

func TestCafeNotOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		url            string
		statusCodeWant int
		bodyWant       string
	}{{url: "/cafe", statusCodeWant: http.StatusBadRequest, bodyWant: "unknown city"},
		{url: "/cafe?city=omsk", statusCodeWant: http.StatusBadRequest, bodyWant: "unknown city"},
		{url: "/cafe?city=tula&count=na", statusCodeWant: http.StatusBadRequest, bodyWant: "incorrect count"},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.url, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, v.statusCodeWant, response.Code)
		assert.Equal(t, v.bodyWant, strings.TrimSpace(response.Body.String()))
	}
}
