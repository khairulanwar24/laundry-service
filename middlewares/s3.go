package middleware

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// FileUploadToS3Middleware middleware for uploading file to S3 with validation
func FileUploadToS3Middleware(c *fiber.Ctx, formFieldName, folderName string, allowedTypes []string, maxFileSize int64) error {
	// Load environment variables from .env file
	err := godotenv.Load("/app/.env")
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}

	// Fetch AWS credentials and S3 bucket info from environment
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("AWS_S3_ENDPOINT")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	// Get the file from form field
	file, err := c.FormFile(formFieldName)
	if err != nil {
		return fmt.Errorf("file upload error: %v", err)
	}

	// Validate file type
	if !isValidFileType(file, allowedTypes) {
		return fmt.Errorf("file type not allowed")
	}

	// Validate file size
	if file.Size > maxFileSize {
		return fmt.Errorf("file size exceeds the maximum limit of %d bytes", maxFileSize)
	}

	// Generate a unique file name
	// fileName := fmt.Sprintf("%s/%d%s", folderName, time.Now().Unix(), filepath.Ext(file.Filename))
	originalFilename := file.Filename
	extension := filepath.Ext(originalFilename) // Dapatkan ekstensi file
	randomFileName := uuid.New().String()       // UUID sebagai nama file
	fileName := fmt.Sprintf("%s/%s_%d%s", folderName, randomFileName, time.Now().Unix(), extension)

	// Upload file to S3
	err = uploadFileToS3(accessKey, secretKey, endpoint, bucketName, fileName, file)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Store the uploaded file name in locals for access in the next handlers
	c.Locals("fileName", fileName)

	return nil
}

// // isValidFileType validates the file type based on the allowed types
// func isValidFileType(file *multipart.FileHeader, allowedTypes []string) bool {
// 	ext := filepath.Ext(file.Filename)
// 	for _, t := range allowedTypes {
// 		if ext == t {
// 			return true
// 		}
// 	}
// 	return false
// }

// uploadFileToS3 uploads the file to the S3 bucket
func uploadFileToS3(accessKey, secretKey, endpoint, bucketName, key string, file *multipart.FileHeader) error {
	// Configure AWS SDK with custom endpoint (IDCloudHost S3)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint, // S3 Endpoint
				SigningRegion: "us-west-2",
			}, nil
		})),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to open file %q, %v", file.Filename, err)
	}
	defer src.Close()

	// Detect content type based on the file extension
	contentType := detectContentType(file.Filename)

	// Upload the file to the S3 bucket
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead, // Set the file to be publicly accessible
	})
	if err != nil {
		return fmt.Errorf("unable to upload file to S3, %v", err)
	}

	return nil
}

// detectContentType returns the content type based on the file extension
func detectContentType(fileName string) string {
	switch filepath.Ext(fileName) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
