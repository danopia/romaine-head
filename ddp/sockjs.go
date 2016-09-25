package ddp

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SockJsInfo struct {
	Websocket    bool
	Origins      []string
	CookieNeeded bool
	Entropy      int
}

func ServeSockJsInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("cache-control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Add("content-type", "application/json; charset=UTF-8")
	w.Header().Add("vary", "origin")

	log.Printf("cb=%s", r.URL.Query().Get("cb"))
	info := SockJsInfo{
		Websocket:    true,
		Origins:      []string{"*:*"},
		CookieNeeded: false,
		Entropy:      int(time.Now().Unix()),
	}

	payload, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(payload)
}
