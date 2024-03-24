package api

import "net/http"

func sendResponse(w http.ResponseWriter, error int) {
  switch error {
  case 200:
    w.WriteHeader(http.StatusOK)
  case 500:
    w.WriteHeader(http.StatusInternalServerError)
  }
}
