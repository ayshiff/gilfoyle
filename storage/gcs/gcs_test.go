package gcs

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	gstorage "cloud.google.com/go/storage"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	assertTest "github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

// NewMockedStorage returns a new mocked Storage.
func NewMockedStorage(ctx context.Context, bucket string, port string) (*Storage, error) {
	client, err := gstorage.NewClient(ctx, option.WithoutAuthentication(), option.WithEndpoint("https://www.googleapis.com/storage/v1/"))
	if err != nil {
		return nil, err
	}
	return &Storage{bucket: client.Bucket(bucket)}, nil
}
func Test(t *testing.T) {
	assert := assertTest.New(t)
	bucket := os.Getenv("GCP_BUCKET")
	port := "4443"

	if bucket == "" {
		bucket = "test-bucket"
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	options := &dockertest.RunOptions{
		Repository: "fsouza/fake-gcs-server",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"4443/tcp": {{HostPort: "4443"}},
		},
	}

	// Build and run the given Dockerfile
	resource, err := pool.RunWithOptions(options, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	endpoint := fmt.Sprintf("localhost:%s", resource.GetPort("4443/tcp"))

	fmt.Printf(endpoint)

	if err = pool.Retry(func() error {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.Get("https://0.0.0.0:4443/storage/v1/b")
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not OK")
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	ctx := context.Background()
	s, err := NewStorage(ctx, bucket, port)
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
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
