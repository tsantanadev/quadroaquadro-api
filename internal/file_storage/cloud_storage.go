package filestorage

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

type FileStorageGCP struct {
	Config *FileStorageConfig
}

type FileStorageConfig struct {
	Bucket string
	Key    string
}

func (s *FileStorageGCP) UploadFile(ctx context.Context, file io.Reader, fileName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create cloud storage client: %w", err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			log.Printf("failed to close cloud storage client: %v", cerr)
		}
	}()

	bkt := client.Bucket(s.Config.Bucket)
	writer := bkt.Object(fileName).NewWriter(ctx)
	writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("failed to upload file to cloud storage: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	return nil
}
