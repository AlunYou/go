package test_util

import (
	"fmt"
	"net/url"
	"net/http"
	"net/http/httptest"
)

func MockHttpServer(code int, body string) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}
	return server, httpClient
}

func CompareSlice(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i:=0; i<len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}