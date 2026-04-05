package repository

import (
	c "github.com/SemgaTeam/semga-stream/internal/config"
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"context"
	"time"
	"net/http"
)

type StorageRepository struct {
	conf *c.Config
	awsConf *aws.Config
	client *s3.Client
	presignClient *s3.PresignClient
}

func NewStorageRepository(conf *c.Config) (*StorageRepository, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(conf.Storage.Region),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.Storage.AccessKeyID, conf.Storage.SecretAccessKey, "")),
        config.WithHTTPClient(&http.Client{}),
  )	
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(conf.Storage.URL)
	})

	presignClient := s3.NewPresignClient(client)

	return &StorageRepository{
		conf,	
		&cfg,
		client,
		presignClient,
	}, nil
}

func (r *StorageRepository) GenerateUploadURL(ctx context.Context, path string) (string, error) {
	req, err := r.presignClient.PresignPutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(r.conf.Storage.Bucket),
			Key:    aws.String(path),
		},
		s3.WithPresignExpires(
			time.Duration(r.conf.Storage.PresignExpirationSeconds) * time.Second,
		),
	)
	if err != nil {
		return "", e.Unknown(err)
	}

	return req.URL, nil
}

func (r *StorageRepository) Delete(ctx context.Context, path string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.conf.Storage.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return e.Unknown(err)
	}

	return nil
}
