package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	var s http.Server
	http2.VerboseLogs = true
	s.Addr = ":8080"

	http2.ConfigureServer(&s, nil)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Hello World")
	})

	log.Fatal(s.ListenAndServeTLS("server.crt", "server.key"))
}
