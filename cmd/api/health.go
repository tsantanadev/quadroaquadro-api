package main

import (
	"net/http"
)

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	if err := writeJSON(w, 200, data); err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
	}
}
