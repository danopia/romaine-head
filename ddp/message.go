package ddp

type Message struct {
	Type       string        `json:"msg,omitempty"`
	Version    string        `json:"version,omitempty"`
	Support    []string      `json:"support,omitempty"`
	Id         string        `json:"id,omitempty"`
	Name       string        `json:"name,omitempty"`
	Params     []interface{} `json:"params,omitempty"`
	ServerId   string        `json:"server_id,omitempty"`
	Session    string        `json:"session,omitempty"`
	Collection string        `json:"collection,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	Subs       []string      `json:"subs,omitempty"`
	Error      map[string]interface{} `json:"error,omitempty"`
	Result     interface{} `json:"result,omitempty"`
}
