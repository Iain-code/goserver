package handler

import (
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
