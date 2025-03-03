package main

import (
	"log"
	"net/http"
)

func main() {

	const port = "8080"
	mux := http.NewServeMux() // checks each request for keywords to know how to handle each request

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/", fileServer)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe()) // starts the server and will "listen" for http requests on given ports and return errors

}
