package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"reco-engine/db"

	"github.com/gorilla/sessions"
)

type Login struct {
  Logger *slog.Logger
  Store *sessions.CookieStore
  Db *db.Db
}

type LoginBody struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

func (l *Login) DoLogin(w http.ResponseWriter, r *http.Request) {
  var login LoginBody
  err := json.NewDecoder(r.Body).Decode(&login)
  if err != nil {
    l.Logger.Error(err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  l.Logger.Info("", 
    slog.String("username", login.Username), 
    slog.String("password", login.Password))
  l.SaveSession(w, r, login.Username)
  w.WriteHeader(http.StatusOK)
}

func (l *Login) DoLogout(w http.ResponseWriter, r *http.Request) {
  l.DelSession(w, r)
  w.WriteHeader(http.StatusOK)
}

func (l *Login) SaveSession(w http.ResponseWriter, 
  r *http.Request, userName string) {
  session, _ := l.Store.Get(r, "cookie")
  session.Values["auth"] = true
  session.Values["user"] = userName
  session.Save(r, w) 
}

func (l *Login) DelSession(w http.ResponseWriter, r *http.Request) {
  session, _ := l.Store.Get(r, "cookie")
  session.Values["auth"] = false
  session.Save(r, w) 
}
