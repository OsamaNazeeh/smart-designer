package api

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)


func presignedGet(ctx context.Context,s3Client *s3.Client, objectKey, bucketName string) (*v4.PresignedHTTPRequest, error){
	presignedClient := s3.NewPresignClient(s3Client)
	
	presignedHTTPRequest, err := presignedClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key: aws.String(objectKey),
		}, 
		s3.WithPresignExpires(15 * time.Minute), 
	)
	if err != nil {
		return nil, err
	}
	return presignedHTTPRequest, nil 
}