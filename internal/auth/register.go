package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/osamaNazieh/smartDesinger/internal/database"
	"github.com/osamaNazieh/smartDesinger/internal/responses"
)

func (h *Handler) SingnUp(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var body Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "Not valid body shape", err)
		return
	}

	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "coudln't hash the password", err)
		return	
	}

	userId := uuid.New()
	accessToken, err := makeJWT(userId, h.app.TokenSecret, 15 * time.Minute)
	fmt.Println(accessToken)
	if err != nil {
		log.Fatal(err)
	}
	userRecord, err := h.app.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: userId,
		Username: body.Username,
		Password: hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Fatal(err)
	}


	responses.RespondWithJSON(w, http.StatusCreated,
	struct {
		ID uuid.UUID `json:"id"`;
		Username string `json:"username"`;
		CreatedAt string `json:"created_at"`;
		UpdatedAt string `json:"updated_at"`;
		AccessToken string `json:"token"`;
	} {
		ID: userRecord.ID,
		Username: userRecord.Username,
		CreatedAt: userRecord.CreatedAt.String(),
		UpdatedAt: userRecord.UpdatedAt.String(),
		AccessToken: accessToken,
	})
}

