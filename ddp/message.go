package ddp

type Message struct {
	Type       string                 `json:"msg,omitempty"`
	Version    string                 `json:"version,omitempty"`
	Support    []string               `json:"support,omitempty"`
	Id         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"` // what uses this?
	Method     string                 `json:"method,omitempty"`
	Params     []interface{}          `json:"params,omitempty"`
	ServerId   string                 `json:"server_id,omitempty"`
	Session    string                 `json:"session,omitempty"`
	Collection string                 `json:"collection,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	Subs       []string               `json:"subs,omitempty"`
	Methods    []string               `json:"methods,omitempty"`
	Error      *ClientError           `json:"error,omitempty"`
	Result     interface{}            `json:"result,omitempty"`
}

type ClientError struct {
	Code    int    `json:"error,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
	Type    string `json:"errorType,omitempty"`
}
