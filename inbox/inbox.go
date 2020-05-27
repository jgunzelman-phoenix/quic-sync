package inbox

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// TopicConfig is a structure that describs a topic you can post too (Name) and the LocalTopic kafka topic
type TopicConfig struct {
	Name       string `json:"name,omitempty"`
	LocalTopic string `json:"local-topic,omitempty"`
}

var availaleKafkaTopics = []string{}
var topicConfigs = make(map[string]*TopicConfig)

// PutMessage on a specified topic
func PutMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topic"]
	topicConfig := topicConfigs[topicName]
	//Handle invalid topic name
	if topicConfig == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Invalid Topic Name. For list of valid topics call api/v1/inbox/topic\"}"))
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	} else {
		w.WriteHeader(http.StatusOK)
		//TODO Add kafka code
	}
}

// GetTopics returns a list of topics available to pot to.
func GetTopics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	topics := []string{}
	for name := range topicConfigs {
		topics = append(topics, name)
	}
	json.NewEncoder(w).Encode(topics)
}

// AddTopicConfig add configuration to quic-sync
func AddTopicConfig(w http.ResponseWriter, r *http.Request) {
	var newTC TopicConfig
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&newTC)
	//failure to unmarhal post.
	if err == nil {
		w.Write([]byte("{\"error\":\"Request body was not formated correctly\"}"))
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		topicConfigs[newTC.Name] = &newTC
		w.WriteHeader(http.StatusOK)
	}
}

// GetTopicConfig return a topic config that corresponds to the provided topic
func GetTopicConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topic"]
	topicConfig := topicConfigs[topicName]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if topicConfig == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Invalid Topic Name. For list of valid topics call api/v1/inbox/topic\"}"))
	} else {
		json.NewEncoder(w).Encode(topicConfig)
	}
}

// GetTopicConfigs return a list of topic configurations
func GetTopicConfigs(w http.ResponseWriter, r *http.Request) {
	topicConfigList := []TopicConfig{}
	for _, value := range topicConfigs {
		topicConfigList = append(topicConfigList, *value)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topicConfigList)
}

// DeleteTopicConfig remove a topic configuration
func DeleteTopicConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topic"]
	//if topic exsists
	if _, ok := topicConfigs[topicName]; ok {
		//TODO delete function
		delete(topicConfigs, topicName)
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Invalid Topic Name. For list of valid topics call api/v1/inbox/topic\"}"))
	}
}
