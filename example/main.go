package main

import (
	"net/http"
	_ "github.com/ezeev/saga/auth0"
	_ "github.com/ezeev/saga/stripe"
	_ "github.com/ezeev/saga/metrics"
	"github.com/ezeev/saga/page"
	"github.com/ezeev/saga/middleware"
	"encoding/json"
)

func init() {
	http.HandleFunc("/", handler)

	// example exposing data needed to render a page via JSON API
	// the middleware package provides rate limiting and sets CORS headers
	// so it can be queried via client Javascript. Ideal for SPAs.
	// If no valid "X-Auth-Token" header is set the page response will only
	// be populated with public-safe info
	http.HandleFunc("/api/page",
		middleware.Api(
		middleware.ApiRateLimit(
			page.HandlePageApi)))


	// example chaining the Api, ApiAuth and ApiRateLimit middleware on an API end point
	// exposing fictional secure data. A user MUST pass a valid JWT token via setting the
	// "X-Auth-Token" header to get a valid response. To see your JWT token, login
	// via the example page, then look for the cn-profile cookie stored in your browser
	// example curl command (replace YOUR_JWT):
	// curl -H "X-Auth-Token: YOUR_JWT" http://localhost:8080/api/example
	http.HandleFunc("/api/example",
		middleware.Api(
		middleware.ApiAuth(
		middleware.ApiRateLimit(
			secureApiExampleHandler))))
}

func handler(w http.ResponseWriter, r *http.Request) {
	page.HandlePage(w,r,"index.html","Home")
}


type mySecureApiData struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

func secureApiExampleHandler(w http.ResponseWriter, r *http.Request) {
	respData := mySecureApiData{Foo:"value 1",Bar:"value 2"}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(respData); err != nil {
		panic(err)
	}
}


