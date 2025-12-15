package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
)

// MoveFile — копирует файл в новую "папку" и удаляет старый
func (c *Client) MoveFile(ctx context.Context, fromKey, toKey string) error {
	_, err := c.client.CopyObject(ctx, &s3sdk.CopyObjectInput{
		Bucket:     aws.String(c.Bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", c.Bucket, fromKey)),
		Key:        aws.String(toKey),
	})
	if err != nil {
		return fmt.Errorf("Ошибка копирования %s → %s: %w", fromKey, toKey, err)
	}

	_, err = c.client.DeleteObject(ctx, &s3sdk.DeleteObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(fromKey),
	})
	if err != nil {
		return fmt.Errorf("Ошибка удаления оригинала %s: %w", fromKey, err)
	}

	return nil
}
