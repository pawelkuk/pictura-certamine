package s3_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/pawelkuk/pictura-certamine/pkg/sdk/s3"
)

func TestMinioClient(t *testing.T) {
	client, err := s3.NewMinioClient(
		os.Getenv("S3_ACCESS_KEY"),
		os.Getenv("S3_SECRET_KEY"),
		os.Getenv("S3_ENDPOINT"),
	)
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}
	ctx := context.Background()
	fileName := fmt.Sprintf("file_%d.txt", time.Now().Unix())
	testContent := "test"
	err = client.PutObject(ctx, fileName, []byte(testContent))
	if err != nil {
		log.Fatalf("could not put object: %v", err)
	}
	buf, err := client.GetObject(ctx, fileName)
	if err != nil {
		log.Fatalf("could not get object: %v", err)
	}
	if string(buf) != testContent {
		log.Fatalf("objects don't match: '%s' != '%s'", buf, testContent)
	}
	list, err := client.ListObjects(ctx)
	if err != nil {
		log.Fatalf("could not list objects: %v", err)
	}
	if !slices.Contains(list, fileName) {
		log.Fatalf("list does not contain %s", fileName)
	}
}
