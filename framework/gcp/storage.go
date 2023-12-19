package gcp

import (
	"context"
	"encoder/application/services"
	"encoder/domain"

	"cloud.google.com/go/storage"
)

type CloudStorageReader struct {
	Bucket  *storage.BucketHandle
	Context context.Context
}

func NewCloudStorageReader(bucketName string) (*CloudStorageReader, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &CloudStorageReader{
		Bucket:  client.Bucket(bucketName),
		Context: ctx,
	}, nil
}

func (c *CloudStorageReader) ReadVideo(video *domain.Video) (services.VideoStorageReader, error) {
	object := c.Bucket.Object(video.FilePath)

	r, err := object.NewReader(c.Context)
	if err != nil {
		return nil, err
	}

	return r, nil
}
