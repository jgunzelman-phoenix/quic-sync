package restApi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	subManager "github.com/jgunzelman-phoenix/quic-sync/qs-kafka"
	model "github.com/jgunzelman-phoenix/quic-sync/qs-model"
)

func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := subManager.DeleteSubscription(id)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	response, _ := json.Marshal(subManager.GetSubscription(id))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(subManager.Subscriptions)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func PutSubscription(w http.ResponseWriter, r *http.Request) {
	newSub := model.Subscription{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newSub)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
		subManager.PutSubscription(&newSub)
	} else {
		errmsg := "{'error':'" + err.Error() + "'}"
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errmsg))
	}
}
