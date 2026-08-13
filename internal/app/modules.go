package app

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/osamaNazieh/smartDesinger/internal/database"
)

// App holds the application dependencies and configuration.
// It contains database access, cloud storage client, and authentication settings.
//
// Fields:
//   - DB: Database query operations handler
//   - TokenSecret: Secret key used for token generation and validation
//   - S3Client: AWS S3 client for object storage operations
//   - BucketName: Name of the S3 bucket for file storage
//   - BucketRegion: AWS region where the S3 bucket is located
type App struct {
	DB           *database.Queries
	TokenSecret  string
	S3Client     *s3.Client
	BucketName   string
	BucketRegion string
}
