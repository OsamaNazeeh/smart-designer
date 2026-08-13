package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/osamaNazieh/smartDesinger/internal/auth"
	"github.com/osamaNazieh/smartDesinger/internal/database"
	"github.com/osamaNazieh/smartDesinger/internal/responses"
)

func (h *Handler )UploadImages(w http.ResponseWriter, r *http.Request) {
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




	fmt.Println("Uploading the image...")
	

	maxMemo := 10 << 20
	if err := r.ParseMultipartForm(int64(maxMemo)); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "Bad Request body", err)
		return 
	} 
	
	file, header, err := r.FormFile("image")
	if err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "Bad Request body", err)
		return 
	}

	defer file.Close()
	key := make([]byte, 16)
	_, err = rand.Read(key)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, "something went wronge at creating object key", err)
		return 
	}
	objectKey := base64.RawURLEncoding.EncodeToString(key)

	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil || (mediaType != "image/jpeg" && mediaType != "image/png") {
		responses.RespondWithError(w, http.StatusBadRequest, "png and jpeg images only allowed", err)
		return
	}

	var ext string 
	switch mediaType {
	case "image/png":
		ext = ".png"
	case "image/jpg":
		ext = ".jpg"
	} 

	if _, err := h.app.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket: aws.String(h.app.BucketName),
		Key: aws.String(objectKey),
		Body: file,
	}); err != nil {
		fmt.Printf("Something went wronge when trying to upload the object, See: %v\n", err)
		return
	}




	record, err := h.app.DB.CreateNewImage(r.Context(), database.CreateNewImageParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Objectkey: objectKey,
		Ext: ext,
		OwnerID: userId,
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	
	
	fmt.Println(record)

	responses.RespondWithJSON(w, http.StatusOK,
		struct {
			ID string `json:"id"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
			UserId string `json:"user_id"`
			URL string `json:"url"`
		} {
			ID: record.ID.String(),
			CreatedAt: record.CreatedAt.Format(time.RFC3339),
			UpdatedAt: record.UpdatedAt.Format(time.RFC3339),
			UserId: record.OwnerID.String(),
			URL: fmt.Sprintf("http://localhost:8080/images/%v", record.ID),
		},
	)	
}