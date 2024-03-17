package main

import (
	"log/slog"
	"net/http"
	"os"

  "reco-engine/api"
  
  "github.com/gorilla/sessions"
)

func main() {
  logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
  logger.Info("Hello, From Recommendation Engine ...")
  defer logger.Info("Bye !!!!!!!")

  store := sessions.NewCookieStore([]byte("Secret-Key"))

  l := api.Login{ Logger: logger }
  m := api.Model{ Logger: logger }

  routes := map[string]api.Adapter{}
  routes["POST /login"] = l.DoLogin
  routes["POST /model"] = m.AddModel

  mux := http.NewServeMux()
  for route, handler := range routes {
    mux.HandleFunc(route, handler)
  }

  wrappedMux := api.NewApiLogger(
    api.NewAuthenticator(mux, logger, store), logger)

  http.ListenAndServe(":8000", wrappedMux)
}
