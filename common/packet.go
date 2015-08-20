package common

// Represents data messages between a head and leaves
type Packet struct {
	Cmd     string
	Context string
	Extras  map[string]interface{}
}
