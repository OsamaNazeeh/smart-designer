package api

import (
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/osamaNazieh/smartDesinger/internal/auth"
	"github.com/osamaNazieh/smartDesinger/internal/responses"
)

func (h *Handler)GetObject(w http.ResponseWriter, r *http.Request) {
	// Get the image id from the url 
	stringifiedImageId := r.PathValue("id")
	
	// Authenticate User
	token, err := auth.GetBearer(r.Header)
	if err != nil || token == "" {
		responses.RespondWithError(w, http.StatusUnauthorized, "no authorization token were provided", err)
		return 
	}
	
	userId, err := auth.ValidateJWT(token, h.app.TokenSecret)
	if err != nil {
		responses.RespondWithError(w, http.StatusUnauthorized, "no authorization token were provided", err)
		return 
	}

	imageID, err := uuid.Parse(stringifiedImageId)
	if err != nil {
		log.Fatal(err)
	}

	record, err := h.app.DB.GetImage(r.Context(), imageID)
	if err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "invalid image id", err)
		return 
	}

	// Check for autherization
	if userId != record.OwnerID {
		responses.RespondWithError(w, http.StatusUnauthorized, "You are not the owner of the image", err)
		return 
	}

	s3Object, err := h.app.S3Client.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: &h.app.BucketName,
		Key: &record.Objectkey,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer s3Object.Body.Close()

	// Important
	if s3Object.ContentType != nil {
		w.Header().Set("Content-Type", *s3Object.ContentType)
	}
	io.Copy(w, s3Object.Body)
}