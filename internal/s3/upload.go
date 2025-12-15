package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadFile — загружает файл в S3
func (c *Client) UploadFile(ctx context.Context, key string, data []byte) error {
	_, err := c.client.PutObject(ctx, &s3sdk.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("Ошибка загрузки файла %s в S3: %w", key, err)
	}
	return nil
}
