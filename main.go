package main

import (
	"database/sql"
	"goserver/handler"
	"goserver/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load() // load env file
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	const port = "8080"
	mux := http.NewServeMux() // checks each request for keywords to know how to handle each request
	apiCfg := &handler.ApiConfig{}
	apiCfg.Db = dbQueries
	apiCfg.Platform = os.Getenv("PLATFORM")
	apiCfg.TokenSecret = os.Getenv("TOKEN_SECRET")
	apiCfg.PolkaKey = os.Getenv("POLKA_KEY")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fileServer := http.FileServer(http.Dir("."))               // http.Dir makes the "." into a filetype that FileServer can read
	wrappedFileServer := http.StripPrefix("/app/", fileServer) // removes the app prefix from the URL then passes the request to fileServer
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(wrappedFileServer))
	mux.HandleFunc("GET /api/healthz", handler.ServerReady)
	mux.HandleFunc("POST /api/chirps", apiCfg.Chirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetOneChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.Counter)
	mux.HandleFunc("POST /admin/reset", apiCfg.Reset)
	mux.HandleFunc("POST /api/users", apiCfg.NewUser)
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("POST /api/refresh", apiCfg.Refresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.Revoke)
	mux.HandleFunc("PUT /api/users", apiCfg.UpdateDetails)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.AddChirpyRed)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe()) // starts the server and will "listen" for http requests on given ports and return errors

}
