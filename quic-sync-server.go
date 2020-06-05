package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/op/go-logging"
	"phoenix-opsgroup.com/quic-sync/inbox"
)

//Constant
const defaultWebPort = 8443
const defaultCertLocation = "./default-certs/server.crt"
const defaultKeyLocation = "./default-certs/server.key"

//Version
const major = 0
const minor = 0
const patch = 0

//Service Variables
var webPort int
var certFile string
var keyFile string
var kafkaBootstrap string

//Logging Variables
var logLevel string
var log = logging.MustGetLogger("example")
var logformat = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	//Set up logging
	flag.StringVar(&logLevel, "log-level", "DEBUG", "Level of Logging")
	log := logging.MustGetLogger("quic-sync-server")
	format := logging.MustStringFormatter("%{color}%{level:.5s}%{color:reset} %{message}")
	logging.SetFormatter(format)
	if logLevel == "DEBUG" {
		logging.SetLevel(logging.DEBUG, "")
	} else if logLevel == "INFO" {
		logging.SetLevel(logging.INFO, "")
	} else if logLevel == "WARN" {
		logging.SetLevel(logging.WARNING, "")
	} else if logLevel == "ERROR" {
		logging.SetLevel(logging.ERROR, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
	//initialize
	log.Info("--- Quic Sync Server ---")
	flag.IntVar(&webPort, "web-port", defaultWebPort, "port to bind to for web server")
	flag.StringVar(&certFile, "cert-file", defaultCertLocation, "cert file for tls")
	flag.StringVar(&keyFile, "key-file", defaultKeyLocation, "key file for tls")
	flag.StringVar(&kafkaBootstrap, "kafka-bootstrap", "localhost:9092", "kafka bootstrap server list ex: host1:port,host2:port")
	flag.Parse()

	log.Debug("CONFIG:")
	log.Debug("web-port: " + strconv.Itoa(webPort))
	log.Debug("cert-file: " + certFile)
	log.Debug("key-file: " + keyFile)
	log.Debug("kafka-bootstrap: " + kafkaBootstrap)

	//Gorilla router initialization
	router := mux.NewRouter()
	log.Info("initalizing web routes")
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
	log.Info("Starting web server ...")
	//Start webserver
	webServerError := http3.ListenAndServeQUIC("localhost:"+strconv.Itoa(webPort), certFile, keyFile, router)
	if webServerError != nil {
		log.Fatal(webServerError)
	}
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]byte("{\"version\":\"" + strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch) + "\"}"))
}
