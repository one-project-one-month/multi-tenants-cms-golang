package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type S3Service struct {
	client *s3.Client
	bucket string
	logger *logrus.Logger
}

func NewS3Service(bucket string, logger *logrus.Logger) (*S3Service, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	svc := &S3Service{
		client: client,
		bucket: bucket,
		logger: logger,
	}
	return svc, nil
}

func (s *S3Service) UploadFile(file *multipart.FileHeader) (*string, error) {
	src, err := file.Open()
	if err != nil {
		s.logger.Errorf("failed to open file: %v", err.Error())
		return nil, err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			s.logger.Errorf("failed to close file: %v", err.Error())
		}
	}(src)
	ext := filepath.Ext(file.Filename)
	filename := "logos/" + time.Now().Format("20060102150405") + ext

	_, err = s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   src,
		ACL:    "public-read",
	})

	if err != nil {
		s.logger.Errorf("failed to upload file: %v", err.Error())
		return nil, err
	}
	fileUrl := fmt.Sprintf("https://%s%s%s", s.bucket, ".s3.amazonaws.com/", filename)
	return &fileUrl, nil

}
