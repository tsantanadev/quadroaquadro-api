package filestorage

import "io"

type FileStorage interface {
	UploadFile(file io.Reader, fileName string) error
}
