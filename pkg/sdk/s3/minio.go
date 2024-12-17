package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/logger"
)

type MinioClient struct {
	client *minio.Client
	log    *logger.Logger
	bucket string
}

func (c *MinioClient) GetObject(ctx context.Context, objName string) ([]byte, error) {
	reader, err := c.client.GetObject(
		context.Background(),
		c.bucket,
		objName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("could not get object: %w", err)
	}
	stat, err := reader.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not get object stat: %w", err)
	}
	// runtime.Breakpoint()
	buff := make([]byte, stat.Size)
	n, err := reader.Read(buff)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("could not read object: %w", err)
	}
	if n != int(stat.Size) {
		return nil, fmt.Errorf("buffer read bytes mismatch %d != %d", n, stat.Size)
	}
	return buff, nil
}
func (c *MinioClient) PutObject(ctx context.Context, objName string, obj []byte) error {
	buff := &bytes.Buffer{}
	n, err := buff.Write(obj)
	if err != nil {
		return fmt.Errorf("could not write to buffer: %w", err)
	}
	if n != len(obj) {
		return fmt.Errorf("write buffer and bytes mismatch %d != %d", n, len(obj))
	}
	info, err := c.client.PutObject(ctx, c.bucket, objName, buff, int64(len(obj)), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("could not put object: %w", err)
	}
	c.log.Info(ctx, "uploaded object", "name", objName, "size", info.Size, "location", info.Location)
	return nil
}
func (c *MinioClient) ListObjects(ctx context.Context) ([]string, error) {
	opts := minio.ListObjectsOptions{Recursive: true}
	objects := []string{}
	for object := range c.client.ListObjects(ctx, c.bucket, opts) {
		if object.Err != nil {
			return nil, fmt.Errorf("could not list objects: %w", object.Err)
		}
		objects = append(objects, object.Key)
	}
	return objects, nil
}

func NewMinioClient(idKey, secretKey, endpoint string) (*MinioClient, error) {
	client, err := minio.New(
		endpoint, // TODO change to some other object storage
		&minio.Options{
			Creds:  credentials.NewStaticV4(idKey, secretKey, ""),
			Secure: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not create minio client: %w", err)
	}
	return &MinioClient{
		client: client,
		log:    logger.New(os.Stdout, logger.LevelInfo, "S3"),
		bucket: "pictura-certamine",
	}, nil
}
