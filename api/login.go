package api

import (
	"log/slog"
	"net/http"
)

type Login struct {
  username string
  password string
  Logger *slog.Logger
}

func (l *Login) DoLogin(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Header().Set("my-status-code", "200")
}
