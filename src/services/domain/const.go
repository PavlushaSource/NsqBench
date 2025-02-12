package domain

type Topic string
type Channel string

const (
	RequestTopic    Topic   = "RequestTopic"
	ResponseTopic   Topic   = "ResponseTopic"
	RequestChannel  Channel = "RequestChannel"
	ResponseChannel Channel = "ResponseChannel"
)
