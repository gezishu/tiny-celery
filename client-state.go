package tinycelery

type clientState uint8

const (
	clientINIT clientState = iota
	clientRUNNING
	clientSTOPPED
)
