package main

import (
	"fmt"
	"log"
	"net/http"
)

func httpLogger(targetMux http.Handler) http.Handler {
	/*
		Like a python decorator

		each request even to the unexistend endpoint have to be logged.
	*/
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		targetMux.ServeHTTP(res, req)

		// My custom logger for each request
		log.Printf("HTTP %s: %s%s from: %s", req.Method, req.Host, req.URL, req.RemoteAddr)
	})
}

func serveMetrics(file_path string, port string, read_token string) {
	log.Println("HTTP endpoint is available on port", port)

	srv := http.NewServeMux()

	srv.HandleFunc("/_metrics", func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query().Get("token")
		if q == read_token {
			http.ServeFile(res, req, file_path)
		} else {
			fmt.Fprintf(res, "Error")
		}

	})

	srv.HandleFunc("/_health", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(res, "{\"status\": \"OK\"}")
	})

	srv.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Welcome")
	})

	log.Fatal(http.ListenAndServe(port, httpLogger(srv)))
}
