package api

import (
  "log/slog"
  "net/http"
)

type Model struct {
  modelName string
  Logger *slog.Logger
}

func (m *Model) AddModel(w http.ResponseWriter, r *http.Request) {

}
