package s3

import "context"

type Client interface {
	GetObject(ctx context.Context, objName string) ([]byte, error)
	PutObject(ctx context.Context, objName string, obj []byte) error
	ListObjects(ctx context.Context) ([]string, error)
}
