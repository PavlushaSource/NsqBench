package domain

type Topic string
type Channel string

const (
	RequestTopic    Topic   = "RequestTopic"
	ResponseTopic   Topic   = "ResponseTopic"
	ResponseChannel Channel = "ResponseChannel"
)
