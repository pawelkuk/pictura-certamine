package s3

import "context"

type FakeClient struct {
	getErr       error
	getObject    []byte
	putErr       error
	ArgPutObject string
	ArgObject    []byte
	listErr      error
	listObjects  []string
}

func (c *FakeClient) GetObject(ctx context.Context, objName string) ([]byte, error) {
	if c.getErr != nil {
		return nil, c.getErr
	}
	return c.getObject, nil
}
func (c *FakeClient) PutObject(ctx context.Context, objName string, obj []byte) error {
	if c.putErr != nil {
		return c.putErr
	}
	c.ArgPutObject = objName
	c.ArgObject = obj
	return nil
}
func (c *FakeClient) ListObjects(ctx context.Context) ([]string, error) {
	if c.listErr != nil {
		return nil, c.listErr
	}
	return c.listObjects, nil
}
