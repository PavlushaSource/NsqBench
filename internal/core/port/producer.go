package port

type Producer interface {
	// Publish synchronously publishes a message body to the specified topic, returning
	// an error if publish failed
	Publish(topicName string, body []byte) error

	// Stop initiates a graceful stop of the Producer (permanent)
	Stop()
}
