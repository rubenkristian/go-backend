package utils

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rubenkristian/backend/configs"
)

type S3Service struct {
	S3Config *configs.S3Config
	s3Client minio.Client
	ctx      context.Context
}

func InitializeS3Client(s3Config *configs.S3Config) (*S3Service, error) {
	client, err := minio.New(s3Config.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Config.AccessKey, s3Config.SecretKey, ""),
		Secure: s3Config.Ssl,
	})

	if err != nil {
		return nil, err
	}

	return &S3Service{
		S3Config: s3Config,
		s3Client: *client,
		ctx:      context.Background(),
	}, nil
}

func (s3Service *S3Service) UploadFile(file string, contentType string, result chan<- error) {
	defer close(result)

	s3Config := s3Service.S3Config

	_, err := s3Service.s3Client.FPutObject(s3Service.ctx, s3Config.BucketName, s3Config.FileLocation, file, minio.PutObjectOptions{ContentType: contentType})

	result <- err
}

func (s3Service *S3Service) GetRegion() string {
	return s3Service.S3Config.Region
}

func (s3Service *S3Service) GetBucketName() string {
	return s3Service.S3Config.BucketName
}
