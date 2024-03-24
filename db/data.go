package db

import (
  "errors"
  "fmt"
)

type Db struct {
  UserModelsMapping map[string]string // model(key) to user(value) mapping
  Models map[string]Model
  Data *Data
  UserData *UserData
}

func GetDb() Db {
  return Db{UserModelsMapping: map[string]string{}, Models: map[string]Model{}}
}

func (db *Db) SetDataForDb(data *Data) {
  db.Data = data
}

func (db *Db) GetModels(user string) []string {
  if _, ok := db.UserModelsMapping[user]; !ok {
    return []string{}
  }
  size := len(db.Models)
  models := make([]string, 0, size)

  for model := range db.Models {
    models = append(models, model)
  }
  return models
}

// This will insert or update the model
func (db *Db) AddModel(user, name string) {
  db.UserModelsMapping[user] = name
  db.Models[name] = Model {features: map[string]bool{}}
}

// This will delete the model if it is present
func (db *Db) DelModel(user, name string) {
  delete(db.UserModelsMapping, user)
  delete(db.Models, name)
}

// -----------------------------------------------------
type Model struct {
  features map[string]bool
}

func (db *Db) GetFeatures(user, modelName string) ([]string, error) {
  if _, ok := db.UserModelsMapping[user]; 
  if _, found := db.Models[modelName]; !found {
    return nil, errors.New(fmt.Sprintf("Model %v not found", modelName))
  }
  
  size := len(db.Models[modelName].features)
  features := make([]string, 0, size)

  for feature := range db.Models[modelName].features {
    features = append(features, feature)
  }
  return features, nil
}

func (db *Db) AddFeatures(user, modelName string, features []string) error {
  if _, found := db.Models[modelName]; !found {
    return errors.New(fmt.Sprintf("Model %v not found", modelName))
  }
  for _, feature := range features {
    db.Models[modelName].features[feature] = true
  }
  return nil
}

func (db *Db) DelFeature(user, modelName string, features []string) (bool, error) {
  if _, found := db.Models[modelName]; !found {
    return false, errors.New(fmt.Sprintf("Model %v not found", modelName))
  }
  msg := make([]string, 0, len(features))
  for _, feature := range features {
    if _, found := db.Models[modelName].features[feature]; !found {
      msg = append(msg, feature)
    } else {
      delete(db.Models[modelName].features, feature)
    }
  }
  if len(msg) > 0 {
    return true, errors.New(fmt.Sprintf("Features %v not found in model %v", 
      msg, modelName))
  }
  return true, nil
}

// -----------------------------------------------------
type Entry struct {
  Name string
  Id int
  Attr map[string]string
}

type Data struct {
  ModelName string
  Entries []Entry
}

func (db *Db) ValidateData() bool {
  modelName := db.Data.ModelName
  modelFeatures := db.Models[modelName].features
  for _, entry := range db.Data.Entries {
    for key := range entry.Attr {
      if present, found := modelFeatures[key]; !found || !present {
        return false
      }
    }
  }
  return true
}

// -----------------------------------------------------
type UserData struct {
  Username string
  Models map[string]Model
}

func (db *Db) AddUserData() {

}
