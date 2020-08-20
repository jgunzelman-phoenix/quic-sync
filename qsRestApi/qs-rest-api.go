package qsRestApi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jgunzelman-phoenix/quic-sync/qsSubManager"

	"github.com/jgunzelman-phoenix/quic-sync/qsModel"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

var MAJOR_VERSION = 0
var ver = qsModel.Version{Version: "0." + strconv.Itoa(MAJOR_VERSION) + ".0"}
var log = logging.MustGetLogger("api")

func GetVersion(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(ver)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Info(BuildLogResponse(r.RemoteAddr, "version"))
}

func BuildLogResponse(path string, host string) string {
	return host + " requested : /quic-sync/" + path
}

func GetTopics(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(qsSubManager.GetTopics())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := qsSubManager.DeleteSubscription(id)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	response, _ := json.Marshal(qsSubManager.GetSubscription(id))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(qsSubManager.Subscriptions)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func PutSubscription(w http.ResponseWriter, r *http.Request) {
	newSub := qsModel.Subscription{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newSub)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
		qsSubManager.PutSubscription(&newSub)
	} else {
		errmsg := "{'error':'" + err.Error() + "'}"
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errmsg))
	}
}
