ptckage main

import (
	"log/slog"
	"net/http"
	"os"

	"reco-engine/api"
	"reco-engine/db"

	"github.com/gorilla/sessions"
)

func main() {
  logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
  logger.Info("Hello, From Recommendation Engine ...")
  defer logger.Info("Bye !!!!!!!")

  db := db.GetDb()

  store := sessions.NewCookieStore([]byte("Secret-Key"))
  store.Options = &sessions.Options{
    Domain: "localhost",
    Path: "/",
    MaxAge: 3600 * 2,
    HttpOnly: false,
  }

  l := api.Login{ Logger: logger, Db: &db, Store: store }
  m := api.Model{ Logger: logger, Db: &db, Store: store }
  f := api.Feature{ Logger: logger, Db: &db, Store: store  }
  d := api.Data{ Logger: logger, Db: &db, Store: store  }

  routes := map[string]api.Adapter{}
  routes["POST /login"] = l.DoLogin
  routes["POST /logout"] = l.DoLogout
  routes["POST /model"] = m.AddModel
  routes["GET /model"] = m.GetModels
  routes["DELETE /model"] = m.DelModel
  routes["POST /feature"] = f.AddFeatures
  routes["GET /feature"] = f.GetFeatures
  routes["DELETE /feature"] = f.DelFeatures
  routes["POST /data"] = d.FillData

  mux := http.NewServeMux()
  for route, handler := range routes {
    mux.HandleFunc(route, handler)
  }

  wrappedMux := api.NewApiLogger(
    api.NewAuthenticator(mux, logger, store), logger)

  http.ListenAndServe(":8000", wrappedMux)
}
