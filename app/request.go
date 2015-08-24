package app

// Request represents what a client might ask
type Request struct {
	Cmd     string
	Args    []string
	Chroot  string
	Context string
	Extras  map[string]interface{}
}
