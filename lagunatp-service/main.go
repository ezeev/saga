package main

import (
	"net/http"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Laguna Technology Partners")
}
func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)

}
