package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
)

// ListObjects — возвращает список ключей в бакете
func (c *Client) ListObjects(ctx context.Context) ([]string, error) {
	resp, err := c.client.ListObjectsV2(ctx, &s3sdk.ListObjectsV2Input{
		Bucket: aws.String(c.Bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("Ошибка списка объектов в бакете %s: %w", c.Bucket, err)
	}

	var keys []string
	for _, obj := range resp.Contents {
		keys = append(keys, *obj.Key)
	}

	return keys, nil
}
