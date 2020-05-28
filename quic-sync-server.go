package main

import (
	"flag"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucas-clemente/quic-go/http3"
	"phoenix-opsgroup.com/quic-sync/inbox"
)

const defaultWebPort = 4343
const defaultCertLocation = "/opt/quic-sync/security/server.cert"
const defaultKeyLocation = "/opt/quic-sync/security/server.key"

var webPort int
var certFile string
var keyFile string
var kafkaBootstrap string

func main() {
	//initialize
	flag.IntVar(&webPort, "web-port", defaultWebPort, "port to bind to for web server")
	flag.StringVar(&certFile, "cert-file", defaultCertLocation, "cert file for tls")
	flag.StringVar(&keyFile, "key-file", defaultKeyLocation, "key file for tls")
	flag.StringVar(&kafkaBootstrap, "kafka-bootstrap", "", "kafka bootstrap server list ex: host1:port,host2:port")

	//webserver initialization
	router := mux.NewRouter()
	//inbox endpoints
	router.HandleFunc("/api/v1/inbox/{topic}", inbox.PutMessage).Methods("POST")
	router.HandleFunc("/api/v1/inbox/topics", inbox.GetTopics).Methods("GET")
	router.HandleFunc("/api/v1/inbox/config", inbox.AddTopicConfig).Methods("PUT")
	router.HandleFunc("/api/v1/inbox/configs", inbox.GetTopicConfigs).Methods("GET")
	router.HandleFunc("/api/v1/inbox/config/{topic}", inbox.GetTopicConfig).Methods("GET")
	router.HandleFunc("/api/v1/inbox/config/{topic}", inbox.DeleteTopicConfig).Methods("DELETE")

	//outbox endpoints
	router.HandleFunc("/api/v1/outbox/topic-connections", inbox.PutMessage).Methods("GET")
	router.HandleFunc("/api/v1/outbox/config/topic-connection", inbox.GetTopics).Methods("PUT")
	router.HandleFunc("/api/v1/outbox/config/topic-connection/{name}", inbox.AddTopicConfig).Methods("GET")
	router.HandleFunc("/api/v1/outbox/config/topic-connection/{name}", inbox.GetTopicConfigs).Methods("DELETE")

	http3.ListenAndServeQUIC("localhost:"+strconv.Itoa(webPort), "/path/to/cert/chain.pem", "/path/to/privkey.pem", router)
}
