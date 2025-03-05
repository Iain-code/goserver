package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) Counter(w http.ResponseWriter, r *http.Request) {

	count := cfg.fileserverHits.Load()
	responseText := fmt.Sprintf(`
 <html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, count)

	w.Write([]byte(responseText))
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {

	if cfg.Platform != "dev" {
		respondWithError(w, 500, "unable to delete users")
		return
	}
	err := cfg.Db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "unable to delete users")
		return
	}
	response := "all users deleted"
	data, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, 500, "Error marshalling JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)

}
