package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)
// https://drstearns.github.io/tutorials/gomiddleware/
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
// https://betterstack.com/community/guides/logging/logging-in-go/

type Adapter func(http.ResponseWriter, *http.Request)

type ApiLogger struct {
  handler http.Handler
  logger *slog.Logger
}

func (l *ApiLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  start := time.Now()
  l.logger.Info("", slog.String("enter", r.URL.Path), slog.String("method", r.Method))

  l.handler.ServeHTTP(w, r)

  l.logger.Info("", slog.String("exit", r.URL.Path),
    slog.String("endtime", time.Since(start).String()))
}

func NewApiLogger(handler http.Handler, logger *slog.Logger) *ApiLogger {
  return &ApiLogger{ handler, logger }
}

type Authenticator struct {
  handler http.Handler
  logger *slog.Logger
  store *sessions.CookieStore
}

func (a *Authenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // only authenticated user allowed
  if r.URL.Path != "/login" {
    session, _ := a.store.Get(r, "cookie")
    if auth, ok  := session.Values["auth"].(bool); !ok || !auth {
      a.logger.Error("Unauthorized user")
      http.Error(w, "Forbidden Access", http.StatusForbidden)
      return
    }
  }

  // serve Http
  a.handler.ServeHTTP(w, r)
}

func NewAuthenticator(handler http.Handler, 
  logger *slog.Logger, store *sessions.CookieStore) *Authenticator {
  return &Authenticator{ handler, logger, store }
}

func getUser(r *http.Request, store *sessions.CookieStore) string {
  session, _ := store.Get(r, "cookie")
  return session.Values["auth"].(string) 
}
