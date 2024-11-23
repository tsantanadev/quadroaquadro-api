package filestorage

import (
	"context"
	"io"
)

type FileStorage interface {
	UploadFile(ctx context.Context, file io.Reader, fileName string) error
}
