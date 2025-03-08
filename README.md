# goserver
go server with auth

Installation
Inside a Go module run "go get github.com/Iain-code/goserver.git"

	
	"GET /api/healthz" - Checks if server is running
	"POST /api/chirps" - Posts a chirp witht body and user ID
	"GET /api/chirps" - Returns all chirps with AUTHOR and SORT query parameters
	"GET /api/chirps/{chirpID}" - Returns chirp using CHIRP ID
	"GET /admin/metrics" - Returns count of hits
	"POST /admin/reset" - Deletes all users 
	"POST /api/users" - Makes a new user with PASSWORD and EMAIL
	"POST /api/login" - Login user with PASSWORD and EMAIL
	"POST /api/refresh" - Refreshs user token
    "POST /api/revoke" - Revokes user access token
	"PUT /api/users" - Updates user details
    "DELETE /api/chirps/{chirpID}" - Deletes chirp using CHIRP ID
	"POST /api/polka/webhooks" - Adds chirpy red TRUE to users account