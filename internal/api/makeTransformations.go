package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/osamaNazieh/smartDesinger/internal/auth"
	"github.com/osamaNazieh/smartDesinger/internal/database"
	"github.com/osamaNazieh/smartDesinger/internal/responses"
)



func (h *Handler)Transform(w http.ResponseWriter, r *http.Request) {
	// Get the image ID 
	stringifiedImageID := r.PathValue("id")
	imageID, err := uuid.Parse(stringifiedImageID)
	if err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "Not valid image id", err)
		return
	}
	
	// Authenticate the user 
	_, err = auth.GetBearer(r.Header)
	if err != nil {
		responses.RespondWithError(w, http.StatusUnauthorized, "Not valid access token", err)
		return
	}


	record, err := h.app.DB.GetImage(r.Context(), imageID)
	if err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "there is no record with this image id", err )
		return 
	}

	result, err := h.app.S3Client.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: aws.String(h.app.BucketName),
		Key: aws.String(record.Objectkey),
	})

	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			log.Printf("Can't get object %s from bucket %s. No such key exists.\n", record.Objectkey, h.app.BucketName)
			err = noKey
		} else {
			log.Printf("Couldn't get object %v:%v. Here's why: %v\n", h.app.BucketName, record.Objectkey, err)
		}
		return
	}

	// Get the filepath 
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, record.Objectkey + record.Ext) 

	// Create temp file 
	tmpFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Schedual removal and closing 
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	defer result.Body.Close()
	
	_ ,err = io.Copy(tmpFile, result.Body)
	if err != nil {
		log.Fatal(err)
	}

	// =======================
	//  Transformations 
	// =======================
	
	// Get the transformation options
	
	
	type Body struct {
		Transformations Transformations`json:"transformations"`
	} 
	var transformation Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transformation); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "bad body request", err)
		return 
	}
		
	options, outputPath, err := buildFFmpegOptions(filePath, transformation.Transformations, h)
	if err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "", err)
		return 
	}
	
	if err := applyTransform(options...); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	outputFile, err := os.Open(outputPath)

	defer os.Remove(outputPath)
	defer outputFile.Close()
	
	fmt.Println(options)

	if _, err := h.app.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket: &h.app.BucketName,
		Key: &record.Objectkey,
		Body: outputFile,
	}); err != nil {
		responses.RespondWithJSON(w, http.StatusInternalServerError, struct{}{})
		return
	}
	
	ext := strings.Split(outputPath, ".")[1]
	fmt.Printf("New ext is %s\n", ext)

	if err := h.app.DB.UpdateImage(r.Context(), database.UpdateImageParams{
		UpdatedAt: time.Now(),
		Ext: ext,
		ID: record.ID,
	}); err != nil {
		log.Fatal(err)
	}


	responses.RespondWithJSON(w, http.StatusOK, struct{
		URL string `json:"url"`
		ID string `json:"id"`
		UpdatedAt string `json:"updated_at"`
	} {
		URL: fmt.Sprintf("images/%s", record.ID.String()),
		ID: record.ID.String(), 
		UpdatedAt: record.UpdatedAt.Format(time.RFC3339),
	})
}