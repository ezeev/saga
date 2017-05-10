package main

import (
	"net/http"
	_ "github.com/ezeev/saga/auth0"
	_ "github.com/ezeev/saga/stripe"
	_ "github.com/ezeev/saga/metrics"
	"github.com/ezeev/saga/page"
	"github.com/ezeev/saga/middleware"
)

func init() {
	http.HandleFunc("/", handler)

	//example exposing account information over API w/ auth and rate limiting
	http.HandleFunc("/api/account",middleware.ApiAuth(page.HandleAccountApi))

}

func handler(w http.ResponseWriter, r *http.Request) {
	page.HandlePage(w,r,"index.html","Home")
}



