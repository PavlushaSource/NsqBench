package domain

type Channel string

const (
	RequestChannel  Channel = "RequestChannel#ephemeral"
	ResponseChannel Channel = "ResponseChannel#ephemeral"
)
