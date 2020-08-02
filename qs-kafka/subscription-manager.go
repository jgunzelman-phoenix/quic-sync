package subManager

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	model "github.com/jgunzelman-phoenix/quic-sync/qs-model"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/op/go-logging"
	kafka "github.com/segmentio/kafka-go"
)

var log = logging.MustGetLogger("sub-manager")
var Subscriptions = make(map[string]model.Subscription)
var roundTripper *http3.RoundTripper
var brokers []string
var cluster sarama.Consumer

const MB = 1048576
const MAX_SIZE = 50 * MB

//Initialize sets up the kafka client and http3 client for transmission
func Initialize(bootstrap string) {
	//Kafka config
	brokers := strings.Split(bootstrap, ",")
	config := sarama.NewConfig()
	cluster, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal("Failed to connect to Kafka!!")
		log.Fatal(err.Error())
		os.Exit(1)
	}
	cluster.Topics()
	//quic confgis
	var keyLog io.Writer
	f, err := os.Create("./keylog")
	if err != nil {
		log.Error(err.Error())
	}
	defer f.Close()
	keyLog = f
	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Error(err.Error())
	}
	var qconf quic.Config
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			KeyLogWriter:       keyLog,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
}

//PutSubscription adds a subscription to the system
func PutSubscription(subscription *model.Subscription) {
	Subscriptions[subscription.Id] = *subscription
	go startSubscriptionThread(subscription)
}

//GetSubscription returns a subscription assoiciated with the id provided
func GetSubscription(id string) *model.Subscription {
	sub := Subscriptions[id]
	return &sub
}

//DeleteSubscription removes the subscription assoicated with the id provided
func DeleteSubscription(id string) error {
	sub := Subscriptions[id]
	if &sub == nil {
		return errors.New("No subscription found")
	} else {
		delete(Subscriptions, id)
		return nil
	}
}

//GetTopics returns a list of available kafka topics
func GetTopics() []string {
	topics, err := cluster.Topics()
	if err != nil {
		log.Error(err.Error())
	}
	return topics
}

func startSubscriptionThread(sub *model.Subscription) {
	id := sub.Id
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "quic-sync-1",
		Topic:    sub.TopicName,
		MinBytes: 0,
		MaxBytes: MAX_SIZE,
	})
	hclient := &http.Client{
		Transport: roundTripper,
	}
	subscription := Subscriptions[id]
	for &subscription != nil {
		msg, err := consumer.ReadMessage(context.Background())
		if &err == nil {
			encodedMsg := base64.StdEncoding.EncodeToString(msg.Value)
			strBytes := []byte(encodedMsg)
			byteReader := bytes.NewReader(strBytes)
			response, err := hclient.Post(subscription.Endpoint, "application/octet-stream", byteReader)
			if err == nil {
				if response.StatusCode != http.StatusOK {
					log.Error("Message sent but recieved " + response.Status)
				}
			} else {
				log.Error(err.Error())
			}
		}
	}

}
