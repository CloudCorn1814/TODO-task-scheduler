package api

import (
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
}
