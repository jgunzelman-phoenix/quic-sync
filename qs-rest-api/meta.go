package restApi

import (
	"encoding/json"
	"net/http"
	model "github.com/jgunzelman-phoenix/quic-sync/qs-model"
)

var ver = model.Version{Version: "0.1.0"}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(ver)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
