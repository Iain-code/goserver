package handler

import (
	"fmt"
	"goserver/internal/auth"
	"net/http"
	"time"
)

func (ApiCfg *ApiConfig) Refresh(w http.ResponseWriter, r *http.Request) {

	type Token struct {
		Token string `json:"token"`
	}
	tkn, err := auth.GetBearerToken(r.Header)
	fmt.Println(tkn)
	if err != nil {
		respondWithError(w, 401, "error collecting token")
		return
	}
	requestedUser, err := ApiCfg.Db.GetUserFromToken(r.Context(), tkn)
	if err != nil {
		respondWithError(w, 401, "token doesnt exist or has expired")
		return
	}

	token, err := auth.MakeJWT(requestedUser.ID, ApiCfg.TokenSecret, time.Hour*1)

	respToken := Token{}
	respToken.Token = token

	respondWithJSON(w, 200, respToken)

}
