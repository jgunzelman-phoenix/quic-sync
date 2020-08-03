package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	subManager "github.com/jgunzelman-phoenix/quic-sync/qs-kafka"
	restApi "github.com/jgunzelman-phoenix/quic-sync/qs-rest-api"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/op/go-logging"
)

//Constants
const defaultWebPort = 8443
const defaultCertLocation = "./configs/server-certs/server.pem"
const defaultKeyLocation = "./configs/server-certs/server.key"

//Service Variables
var httpsPort int
var http3Port int
var certFile string
var keyFile string
var kafkaBootstrap string

//Logging Variables
var logLevel string
var Log logging.Logger

func main() {
	//Set up logging
	flag.StringVar(&logLevel, "log-level", "DEBUG", "Level of Logging")
	Log := logging.MustGetLogger("quic-sync-server")
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

	//Initialize
	Log.Info("--- Quic Sync Server ---")
	flag.IntVar(&httpsPort, "https-port", defaultWebPort, "port to bind to for https server")
	flag.IntVar(&http3Port, "http3-port", defaultWebPort, "port to bind to for http3 server")
	flag.StringVar(&certFile, "cert-file", defaultCertLocation, "cert file for tls")
	flag.StringVar(&keyFile, "key-file", defaultKeyLocation, "key file for tls")
	flag.StringVar(&kafkaBootstrap, "kafka-bootstrap", "localhost:9092", "kafka bootstrap server list ex: host1:port,host2:port")
	flag.Parse()
	Log.Debug("CONFIG:")
	Log.Debug("https-port      : " + strconv.Itoa(httpsPort))
	Log.Debug("http3-port      : " + strconv.Itoa(http3Port))
	Log.Debug("cert-file       : " + certFile)
	Log.Debug("key-file        : " + keyFile)
	Log.Debug("kafka-bootstrap : " + kafkaBootstrap)

	//Set initialize Subscription Manager
	subManager.Initialize(kafkaBootstrap)

	//Gorilla router initialization
	router := mux.NewRouter()
	Log.Info("initalizing web routes")

	//meta endpoints
	router.HandleFunc("/quic-sync/version", restApi.GetVersion).Methods("GET")

	//sync endpoints
	router.HandleFunc("/quic-sync/v0/message/{topic}", restApi.PostMessage).Methods("POST")
	router.HandleFunc("/quic-sync/v0/topics", restApi.GetTopics).Methods("GET")

	//Subscription endpoints
	router.HandleFunc("/quic-sync/v0/subscriptions", restApi.PutSubscription).Methods("PUT")
	router.HandleFunc("/quic-sync/v0/subscriptions", restApi.GetSubscriptions).Methods("GET")
	router.HandleFunc("/quic-sync/v0/subscription/{id}", restApi.GetSubscription).Methods("GET")
	router.HandleFunc("/quic-sync/v0/subscription/{id}", restApi.DeleteSubscription).Methods("DELETE")

	//Start webserver
	go starthttp3(router)
	starthttps(router)
}

func starthttps(router *mux.Router) {
	Log.Info("Starting https web server ...")
	webServerError := http.ListenAndServeTLS("0.0.0.0:"+strconv.Itoa(httpsPort), certFile, keyFile, router)
	if webServerError != nil {
		Log.Error(webServerError.Error())
	}
}

func starthttp3(router *mux.Router) {
	Log.Info("Starting http3 quic web server ...")
	webServerError := http3.ListenAndServeQUIC("0.0.0.0:"+strconv.Itoa(http3Port), certFile, keyFile, router)
	if webServerError != nil {
		Log.Error(webServerError.Error())
	}
}
