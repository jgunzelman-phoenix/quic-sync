package qsModel

// A topic connection represents a connection to anothe quic sync server.  This config will post message from the local topic to the url provided
type Subscription struct {
	// Unique Id of the connection
	Id string `json:"id,omitempty"`
	// Name off the certificate that requested the connection.
	UserName string `json:"user_name,omitempty"`
	// Designates wheather the subscription is enabled
	Enabled bool `json:"enabled,omitempty"`
	// Name of the local kafka topic that is read for this subscription out to other quic sync servers.
	TopicName string `json:"topic_name,omitempty"`
	// Remote target url that data from the local-topic will be sent to.
	Endpoint string `json:"endpoint,omitempty"`
}

//Version structure
type Version struct {
	Version string `json:"version,omitempty"`
}
