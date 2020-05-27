package outbox

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// TopicConnection : Describes a connection to a kafka topic and a url to post to when a
type TopicConnection struct {
	name       string
	localTopic string
	url        string
}

var topicConnections = make(map[string]*TopicConnection)

// GetTopicConnections : returns a list of topic connections
func GetTopicConnections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topicConnections)
}

// AddTopicConnections : adds Subcription Connection to the system
func AddTopicConnections(w http.ResponseWriter, r *http.Request) {
	var newTC TopicConnection
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&newTC)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//add sub
		topicConnections[newTC.name] = &newTC
	}
}

// GetTopicConnection : get a topic connection correponding to the provided name
func GetTopicConnection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	topicConn := topicConnections[name]
	if topicConn != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(topicConn)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]byte("{\"error\":\"Invalid Topic Connection name\"}"))
	}
}

// DeleteTopicConnection : delete a topic connection correcponding to the provided name
func DeleteTopicConnection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	topicConn := topicConnections[name]
	if topicConn != nil {
		delete(topicConnections, topicConn.name)
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]byte("{\"error\":\"Invalid Topic Connection name\"}"))
	}
}
