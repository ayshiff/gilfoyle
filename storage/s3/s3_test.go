package s3_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage"
	assertTest "github.com/stretchr/testify/assert"
)

func TestS3(t *testing.T) {
	assert := assertTest.New(t)
	port := os.Getenv("FAKE_S3_SERVER_PORT")

	if port == "" {
		port = "7000"
	}

	gilfoyle.Config.Storage.S3 = config.S3Config{
		Hostname:        "127.0.0.1:" + port,
		AccessKeyID:     "access_key",
		SecretAccessKey: "secret_key",
		Bucket:          "gilfoyle-aws-bucket",
		EnableSSL:       false,
	}

	t.Run("should return error file does not exist", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		_, err = s.Stat(ctx, "doesnotexist")
		assert.EqualError(err, "The specified key does not exist.")
	})

	t.Run("should create file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)
	})

	t.Run("should get metadata of file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

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
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		err = s.Delete(ctx, "world")
		assert.NoError(err)

		_, err = s.Stat(ctx, "world")
		assert.EqualError(err, "The specified key does not exist.")
	})

	t.Run("should create then open file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

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

	t.Run("should delete the file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		err = s.Delete(ctx, "world")
		assert.NoError(err)
	})
}
