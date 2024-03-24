package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"reco-engine/db"

	"github.com/gorilla/sessions"
)

type Feature struct {
  Logger *slog.Logger
  Db *db.Db
  Store *sessions.CookieStore
}

type FeatureBody struct {
  ModelName string `json:"model_name"`
  Features []string `json:"features"`
}

func (f *Feature) GetFeatures(w http.ResponseWriter, r *http.Request) {
  // get all features with model
  var feature FeatureBody
  err := json.NewDecoder(r.Body).Decode(&feature)
  if err != nil {
    f.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  user := getUser(r, f.Store)
  features, err := f.Db.GetFeatures(user, feature.ModelName)
  if err != nil {
    f.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Add("Content-Type", "application/json")
  err = json.NewEncoder(w).Encode(features)
  if err != nil {
    f.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func (f *Feature) AddFeatures(w http.ResponseWriter, r *http.Request) {
  // delete features with model
  var feature FeatureBody
  err := json.NewDecoder(r.Body).Decode(&feature)
  if err != nil {
    f.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  err = f.Db.AddFeatures(feature.ModelName, feature.Features)
  if err != nil {
    f.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  logMsg := fmt.Sprintf("Inserted features %v in model %v", 
    feature.Features, feature.ModelName)
  f.Logger.Info(logMsg)
}

func (f *Feature) DelFeatures(w http.ResponseWriter, r *http.Request) {
  // delete features with model
  var feature FeatureBody
  err := json.NewDecoder(r.Body).Decode(&feature)
  if err != nil {
    f.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  ok, err := f.Db.DelFeature(feature.ModelName, feature.Features)
  if !ok && err != nil {
    f.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  logMsg := fmt.Sprintf("Deleted features %v in model %v", 
    feature.Features, feature.ModelName)
  f.Logger.Info(logMsg)
}
