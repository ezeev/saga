package main

import (
	"net/http"
	"github.com/ezeev/saga/auth0"
	"github.com/ezeev/saga/stripe"
	"github.com/ezeev/saga/page"
)

func init() {
	http.HandleFunc("/", handler)

	auth0.RegisterHandlers()
	stripe.RegisterHanlders()


}

func handler(w http.ResponseWriter, r *http.Request) {
	page.HandlePage(w,r,"index.html","Home")
}



