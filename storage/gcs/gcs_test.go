package gcs

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	gstorage "cloud.google.com/go/storage"
	assertTest "github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

// NewMockedStorage returns a new mocked Storage.
func NewMockedStorage(ctx context.Context, bucket string, port string) (*Storage, error) {
	client, err := gstorage.NewClient(ctx, option.WithoutAuthentication(), option.WithEndpoint("http://localhost:"+port))
	if err != nil {
		return nil, err
	}
	return &Storage{bucket: client.Bucket(bucket)}, nil
}

func Test(t *testing.T) {
	assert := assertTest.New(t)
	bucket := os.Getenv("GCP_BUCKET")
	port := os.Getenv("FAKE_GCS_SERVER_PORT")

	if bucket == "" {
		bucket = "test-bucket"
	}

	ctx := context.Background()
	s, err := NewMockedStorage(ctx, bucket, port)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("should return error file does not exist", func(t *testing.T) {
		assert.NoError(err)

		ctx := context.Background()

		_, err = s.Stat(ctx, "doesnotexist")
		assert.EqualError(err, "The specified key does not exist.")
	})

	t.Run("should create file", func(t *testing.T) {
		ctx := context.Background()

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)
	})

	t.Run("should get metadata of file", func(t *testing.T) {
		ctx := context.Background()

		before := time.Now().Add(-1 * time.Second)

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		now := time.Now().Add(2 * time.Second)

		stat, err := s.Stat(ctx, "world")
		assert.NoError(err)

		assert.Equal(int64(5), stat.Size)
		assert.Equal(false, stat.ModifiedTime.Before(before))
		assert.Equal(false, stat.ModifiedTime.After(now))
	})

	t.Run("should create then delete file", func(t *testing.T) {
		ctx := context.Background()

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		err = s.Delete(ctx, "world")
		assert.NoError(err)

		_, err = s.Stat(ctx, "world")
		assert.EqualError(err, "The specified key does not exist.")
	})

	t.Run("should create then open file", func(t *testing.T) {
		ctx := context.Background()

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		f, err := s.Open(ctx, "world")
		assert.NoError(err)
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		assert.NoError(err)
		assert.Equal("hello", string(b))
	})
}
