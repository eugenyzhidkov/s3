package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"

	"s3/internal/config"
	"s3/internal/infrastructure/logger"
)

type Client struct {
	client *s3sdk.Client
	Bucket string
}

func New(cfg *config.Config) (*Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.S3UseSSL {
			return aws.Endpoint{URL: fmt.Sprintf("https://%s", cfg.S3Endpoint)}, nil
		}
		return aws.Endpoint{URL: fmt.Sprintf("http://%s", cfg.S3Endpoint)}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKey,
			cfg.S3SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания AWS конфигурации: %w", err)
	}

	s3Options := func(o *s3sdk.Options) {
		o.UsePathStyle = true
	}

	s3Client := s3sdk.NewFromConfig(awsCfg, s3Options)

	logger.Info("S3-клиент создан", "endpoint", cfg.S3Endpoint, "bucket", cfg.S3Bucket)

	return &Client{
		client: s3Client,
		Bucket: cfg.S3Bucket,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.ListBuckets(ctx, &s3sdk.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к S3: %w", err)
	}
	logger.Info("Подключение к S3 проверено успешно")
	return nil
}
