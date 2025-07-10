package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nats-io/nats.go"
)

type FileUploadMessage struct {
	FileName    string `json:"file_name"`
	FileContent []byte `json:"file_content"`
	BucketName  string `json:"bucket_name"`
	ObjectKey   string `json:"object_key"`
}

var (
	natsConn *nats.Conn
	s3Client *s3.Client
)

func main() {
	// Initialize AWS S3 client
	cfg := awsConfig()
	s3Client = s3.NewFromConfig(*cfg)

	// Initialize NATS connection
	var err error
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	natsConn, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer natsConn.Close()

	// Subscribe to file upload messages
	_, err = natsConn.Subscribe("file.upload", handleFileUpload)
	if err != nil {
		log.Fatal("Failed to subscribe to NATS:", err)
	}

	log.Println("Service B started - listening for file upload messages...")

	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)
	<-osCh

	log.Println("Shutting down...")
}

func handleFileUpload(msg *nats.Msg) {
	var uploadMsg FileUploadMessage

	err := json.Unmarshal(msg.Data, &uploadMsg)
	if err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Received file upload request for: %s", uploadMsg.FileName)

	// Upload to S3
	err = uploadToS3(uploadMsg)
	if err != nil {
		log.Printf("Error uploading to S3: %v", err)
		return
	}

	log.Printf("Successfully uploaded %s to S3", uploadMsg.FileName)
}

func uploadToS3(uploadMsg FileUploadMessage) error {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(uploadMsg.BucketName),
		Key:    aws.String(uploadMsg.ObjectKey),
		Body:   bytes.NewReader(uploadMsg.FileContent),
	})
	return err
}

func awsConfig() *aws.Config {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" || region == "" {
		log.Fatal("AWS credentials are missing")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey, secretKey, "",
			)),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
