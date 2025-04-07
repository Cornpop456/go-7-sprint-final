package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		require.Equal(t, http.StatusOK, response.Code)
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

		require.Equal(t, v.statusCodeWant, response.Code)
		assert.Equal(t, v.bodyWant, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"
	wantedStatus := http.StatusOK

	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: min(len(cafeList[city]), 100)},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&count=%d", city, v.count), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, wantedStatus, response.Code)

		cities := strings.Split(strings.TrimSpace(response.Body.String()), ",")

		if cities[0] == "" {
			cities = cities[1:]
		}

		assert.Equal(t, v.want, len(cities))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"
	wantedStatus := http.StatusOK

	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&search=%s", city, v.search), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, wantedStatus, response.Code)

		cities := strings.Split(strings.TrimSpace(response.Body.String()), ",")

		if cities[0] == "" {
			cities = cities[1:]
		}

		for _, city := range cities {
			assert.True(t, strings.Contains(strings.ToLower(city), strings.ToLower(v.search)))
		}

		assert.Equal(t, v.wantCount, len(cities))
	}
}
