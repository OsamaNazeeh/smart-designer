package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/osamaNazieh/smartDesinger/internal/responses"
)

func GetBearer(header http.Header) (string, error) {
	bearer := header.Get("Authorization")
	if bearer == "" {
		return "", errors.New("No Authorization header were provided")
	}
	token := strings.TrimPrefix(bearer, "Bearer ")
	return token, nil 
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {	
	
	var body Body
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "not valid request body", err)
		return 
	}

	record, err := h.app.DB.GetUserByName(r.Context(), body.Username)
	if err != nil {
		responses.RespondWithError(w, http.StatusUnauthorized, "incorrect credintials", err)
		return
	}

	valid, err := validatePasswordHash(body.Password, record.Password)
	if !valid || err != nil {
		responses.RespondWithError(w, http.StatusUnauthorized, "There is no such user", err)
		return
	}

	accessToken, err := makeJWT(record.ID, h.app.TokenSecret, time.Minute * 15)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, "something went wronge", err)
		return
	}

	responses.RespondWithJSON(w, http.StatusOK, struct{
		ID string `json:"id"`
		Username string `json:"username"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		AccessToken string `json:"token"`
	} {
		ID: record.ID.String(),
		Username: record.Username,
		CreatedAt: record.CreatedAt.String(),
		UpdatedAt: record.UpdatedAt.String(),
		AccessToken: accessToken,
	})
}