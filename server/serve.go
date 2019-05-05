package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lpar/gzipped"
)

func main() {
	fs := withIndices(gzipped.FileServer(http.Dir(".")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", port),
			fs,
		),
	)
}

func withIndices(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = r.URL.Path + "index.html"
		}
		h.ServeHTTP(w, r)
	})
}
