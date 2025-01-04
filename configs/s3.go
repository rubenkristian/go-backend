package configs

import (
	"os"
)

type S3Config struct {
	EndPoint         string
	AccessKey        string
	SecretKey        string
	Region           string
	Ssl              bool
	BucketName       string
	FileLocation     string
	UploadFileMethod string
}

func (env *EnvConfig) LoadS3Config() *S3Config {
	return &S3Config{
		EndPoint:         os.Getenv("S3_ENDPOINT"),
		AccessKey:        os.Getenv("S3_ACCESS_KEY"),
		SecretKey:        os.Getenv("S3_SECRET_KEY"),
		Region:           os.Getenv("S3_REGION"),
		Ssl:              os.Getenv("S3_SSL") == "TRUE",
		BucketName:       os.Getenv("S3_BUCKET_NAME"),
		FileLocation:     os.Getenv("S3_FILE_LOCATION"),
		UploadFileMethod: os.Getenv("UPLOAD_FILE_METHOD"),
	}
}
