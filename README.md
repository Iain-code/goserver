# goserver
go server with auth

Installation
Inside a Go module run "go get github.com/Iain-code/goserver.git"

	
	"GET /api/healthz"
	"POST /api/chirps"
	"GET /api/chirps"
	"GET /api/chirps/{chirpID}"
	"GET /admin/metrics"
	"POST /admin/reset"
	"POST /api/users"
	"POST /api/login"
	"POST /api/refresh"
    "POST /api/revoke"
	"PUT /api/users"
    "DELETE /api/chirps/{chirpID}"
	"POST /api/polka/webhooks"