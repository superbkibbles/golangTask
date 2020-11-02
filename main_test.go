package main

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init(){
	DB , dbErr = ConnectDb("root:{PASSWORD}@/{DBNAME}")
}

func TestCreateEndpoint(t *testing.T) {

	t.Parallel()

	bodyReader := strings.NewReader(`{"longUrl": "https://www.youtube.com"}`)
	req, _ := http.NewRequest("PUT", "/create", bodyReader)
	resp := httptest.NewRecorder()

	Router("/create", "PUT", CreateEndpoint).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "expecting status 200")
}

func TestShowEndpoint(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequest("GET", "/show/?shortUrl=9pK64W3", nil)
	resp := httptest.NewRecorder()

	Router("/show/", "GET", ShowEndpoint).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "expecting status 200")
}

func TestRootEndpoint(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequest("GET", "/rE3ZOOk", nil)
	resp := httptest.NewRecorder()

	Router("/{id}", "GET", RootEndpoint).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusFound, resp.Code)
}

func Router(url string, method string, endpoint http.HandlerFunc) *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc(url, endpoint).Methods(method)

	return router
}