package main

import (
	"github.com/gorilla/mux"
	"github.com/lucas-clemente/quic-go/http3"
	"phoenix-opsgroup.com/quic-sync/inbox"
)

func main() {
	router := mux.NewRouter()
	//inbox endpoints
	router.HandleFunc("/api/v1/inbox/{topic}", inbox.PutMessage).Methods("POST")
	router.HandleFunc("/api/v1/inbox/topics", inbox.GetTopics).Methods("GET")
	router.HandleFunc("/api/v1/inbox/config", inbox.AddTopicConfig).Methods("PUT")
	router.HandleFunc("/api/v1/inbox/configs", inbox.GetTopicConfigs).Methods("GET")
	router.HandleFunc("/api/v1/inbox/config/{topic}", inbox.GetTopicConfig).Methods("GET")
	router.HandleFunc("/api/v1/inbox/config/{topic}", inbox.DeleteTopicConfig).Methods("DELETE")
	//outbox endpoints

	http3.ListenAndServeQUIC("localhost:4242", "/path/to/cert/chain.pem", "/path/to/privkey.pem", router)
}
