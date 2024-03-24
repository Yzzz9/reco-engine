package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"reco-engine/db"

	"github.com/gorilla/sessions"
)

type Data struct {
  Logger *slog.Logger
  Db *db.Db
  Store *sessions.CookieStore
}

type Attr struct {
  Name string `json:"name"`
  Value string `json:"value"`
}

type Entry struct {
  Name string `json:"name"`
  Attr []Attr `json:"attr"`
}

type DataBody struct {
  ModelName string `json:"model_name"`
  Entries []Entry `json:"entries"`
}

func (d *Data) FillData(w http.ResponseWriter, r *http.Request) {
  // make the check against all the features for each entry of data
  // Then insert the data into tables
  var data DataBody
  err := json.NewDecoder(r.Body).Decode(&data)
  if err != nil {
    d.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // Convert the data into db structures
  newData := db.Data{ ModelName: data.ModelName, Entries: []db.Entry{} }
  for i, entry := range data.Entries {
    dbEntry := db.Entry {Name: entry.Name, Id: i, Attr: map[string]string{} }
    for _, element := range entry.Attr {
      dbEntry.Attr[element.Name] = element.Value
    }
    newData.Entries = append(newData.Entries, dbEntry)
  }

  d.Db.SetDataForDb(&newData)

  // validate the requet
  if !d.Db.ValidateData() {
    d.Logger.Error("Validation of data is failed")
    http.Error(w, err.Error(), http.StatusBadRequest)
  } else {
    d.Logger.Info("Validation of data against feature successful")
  }
}
