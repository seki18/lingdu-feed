package common

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appconfig "github.com/seki18/lingdu-feed/config"
)

// S3Client is the shared S3 client singleton.
var (
	S3Client *s3.Client
	S3Bucket string
	S3Region string
)

// InitS3 initializes the AWS S3 client.
func InitS3(cfg appconfig.Config) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Printf("[S3] Failed to load AWS config: %v", err)
		S3Client = nil
		return
	}

	S3Client = s3.NewFromConfig(awsCfg)
	S3Bucket = cfg.S3Bucket
	S3Region = cfg.AWSRegion

	// Quick connectivity check: HeadBucket
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.S3Bucket),
	}); err != nil {
		log.Printf("[S3] Bucket check failed (degraded mode): %v", err)
		S3Client = nil
	} else {
		log.Println("[S3] Connected successfully")
	}
}
