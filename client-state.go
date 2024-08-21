package tinycelery

type ClientState uint8

const (
	ClientINIT ClientState = iota + 1
	ClientRUNNING
	ClientSTOPPED
)
