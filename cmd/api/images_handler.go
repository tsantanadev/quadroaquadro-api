package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const MAX_BODY_SIZE = 10 << 20

func (app *application) imagePostHandler(w http.ResponseWriter, r *http.Request) {
	movieId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(movieId)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	movieExist, err := app.store.Movies.Exists(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !movieExist {
		writeJSONError(w, http.StatusNotFound, "Movie not found")
		return
	}

	http.MaxBytesReader(w, r.Body, MAX_BODY_SIZE)

	if err := r.ParseMultipartForm(MAX_BODY_SIZE); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	for key := range r.MultipartForm.File {
		files := r.MultipartForm.File[key]
		if len(files) == 0 {
			writeJSONError(w, http.StatusBadRequest, "Image is required")
			return
		}

		for _, fileHeader := range files {

			file, err := fileHeader.Open()
			if err != nil {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
				return

			}

			defer file.Close()

			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			filetype := http.DetectContentType(buff)
			if filetype != "image/jpeg" && filetype != "image/png" {
				http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
				return
			}

			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fileName := fmt.Sprintf("%s-%d", key, time.Now().Unix())

			err = app.fileStorage.UploadFile(file, fileName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
