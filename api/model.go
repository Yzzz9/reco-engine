package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"reco-engine/db"

	"github.com/gorilla/sessions"
)

type Model struct {
  Logger *slog.Logger
  Db *db.Db
  Store *sessions.CookieStore
}

type ModelBody struct {
  ModelName string `json:"model_name"`
}

func (m *Model) GetModels(w http.ResponseWriter, r *http.Request) {
  user := getUser(r, m.Store)
  models := m.Db.GetModels(user)

  w.Header().Add("Content-Type", "application/json")
  err := json.NewEncoder(w).Encode(models)
  if err != nil {
    m.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func (m *Model) AddModel(w http.ResponseWriter, r *http.Request) {
  // create model with user
  var model ModelBody
  err := json.NewDecoder(r.Body).Decode(&model)
  if err != nil {
    m.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return 
  }

  user := getUser(r, m.Store)
  m.Db.AddModel(user, model.ModelName)

  logMsg := fmt.Sprintf("Added model %v", model.ModelName)
  m.Logger.Info(logMsg)
}

func (m *Model) DelModel(w http.ResponseWriter, r *http.Request) {
  // delete model with user
  var model ModelBody
  err := json.NewDecoder(r.Body).Decode(&model)
  if err != nil {
    m.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return 
  }

  user := getUser(r, m.Store)
  m.Db.DelModel(user, model.ModelName)

  logMsg := fmt.Sprintf("Deleted model %v", model.ModelName)
  m.Logger.Info(logMsg)
}
