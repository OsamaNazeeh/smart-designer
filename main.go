package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/osamaNazieh/smartDesinger/internal/api"
	"github.com/osamaNazieh/smartDesinger/internal/app"
	"github.com/osamaNazieh/smartDesinger/internal/auth"
	"github.com/osamaNazieh/smartDesinger/internal/database"
)


func connectToDB(dbUrl string) *database.Queries {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	return database.New(db) 
}

func getS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(
		context.Background(), 
		config.WithSharedConfigProfile("smart-designer"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return s3.NewFromConfig(cfg)
}



func main () {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	bucketName := os.Getenv("BUCKET_NAME")
	bucketRegion := os.Getenv("BUCKET_REGION")
	
	queries := connectToDB(dbUrl)
	
	serverHandler := &app.App{
		DB: queries, 
		TokenSecret: tokenSecret,
		S3Client: getS3Client(),
		BucketName: bucketName,
		BucketRegion: bucketRegion,
	}


	

	authHandler := auth.NewHandler(serverHandler)
	apiHandler := api.NewHandler(serverHandler)
	fs := http.FileServer(http.Dir("./static"))

	mux := http.NewServeMux()
	
	mux.Handle("/", fs)


	mux.HandleFunc("POST /register", authHandler.SingnUp)
	mux.HandleFunc("POST /login", authHandler.Login)


	mux.HandleFunc("POST /images", apiHandler.UploadImages)
	mux.HandleFunc("GET /images/{id}", apiHandler.GetObject)
	mux.HandleFunc("POST /images/{id}/transform", apiHandler.Transform)
	

	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}