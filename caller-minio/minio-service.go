package main

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioHelper wraps the MinIO client
type MinioService struct {
	Client     *minio.Client
	BucketName string
}

// New creates and initializes the MinIO client
func NewMinio(endpoint, accessKey, secretKey, bucketName string) (*MinioService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	ctx := context.Background()
	// Ensure bucket exists
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinioService{
		Client:     client,
		BucketName: bucketName,
	}, nil
}

// Upload uploads data to the bucket
func (m *MinioService) Upload(ctx context.Context, objectName string, data []byte, contentType string) error {
	_, err := m.Client.PutObject(ctx, m.BucketName, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Download retrieves an object from the bucket
func (m *MinioService) Download(ctx context.Context, objectName string) ([]byte, error) {
	obj, err := m.Client.GetObject(ctx, m.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	return io.ReadAll(obj)
}
