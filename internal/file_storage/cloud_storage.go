package filestorage

import (
	"context"
	"io"
	"log/slog"

	"cloud.google.com/go/storage"
)

type FileStorageGCP struct {
	Config *FileStorageConfig
}

type FileStorageConfig struct {
	Bucket string
	Key    string
}

func (s *FileStorageGCP) UploadFile(file io.Reader, fileName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		slog.Error("Failed to create cloud storage client: %v", err)
		return err
	}
	defer client.Close()
	bkt := client.Bucket(s.Config.Bucket)

	writer := bkt.Object(fileName).NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		slog.Error("Failed to upload file to cloud storage: %v", err)
		return err
	}

	if err := writer.Close(); err != nil {
		slog.Error("Failed to close writer: %v", err)
		return err
	}

	return nil
}
