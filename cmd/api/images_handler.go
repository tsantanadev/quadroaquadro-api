package main

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	filestorage "github.com/tsantanadev/quadroaquadro/internal/file_storage"
	"github.com/tsantanadev/quadroaquadro/internal/store"
)

const MAX_BODY_SIZE = 10 << 20 // 10 MB

func (app *application) imagePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getMovieIDFromRequest(r)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := checkMovieExists(r.Context(), app.store, id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := parseMultipartForm(w, r); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	files, err := getFilesFromRequest(r)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	for key, fileHeader := range files {
		imageId, err := uploadImage(fileHeader, app.fileStorage, r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		image := store.Image{
			ID:      imageId,
			MovieId: id,
			Level:   key,
		}
		if err := app.store.Images.Create(&image); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Files uploaded successfully"})
}

func getMovieIDFromRequest(r *http.Request) (int, error) {
	movieId := chi.URLParam(r, "id")
	return strconv.Atoi(movieId)
}

func checkMovieExists(ctx context.Context, store store.Storage, id int) error {
	exists, err := store.Movies.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("movie not found")
	}
	return nil
}

func parseMultipartForm(w http.ResponseWriter, r *http.Request) error {
	http.MaxBytesReader(w, r.Body, MAX_BODY_SIZE)
	return r.ParseMultipartForm(MAX_BODY_SIZE)
}

func getFilesFromRequest(r *http.Request) (map[int]*multipart.FileHeader, error) {
	files := make(map[int]*multipart.FileHeader)
	for key, file := range r.MultipartForm.File {
		// return error if key is not a number
		level, err := strconv.Atoi(key)

		if err != nil {
			return nil, fmt.Errorf("invalid level, only numbers are allowed")
		}

		if len(file) == 1 {
			files[level] = file[0]
		} else {
			return nil, fmt.Errorf("only one image per level")
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("image is required")
	}
	return files, nil
}

func uploadImage(fileHeader *multipart.FileHeader, storage filestorage.FileStorage, ctx context.Context) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", err
	}

	// Reset the file pointer to the beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	fileName := uuid.New().String()
	if err := storage.UploadFile(ctx, file, fileName); err != nil {
		return "", err
	}
	return fileName, nil

}
