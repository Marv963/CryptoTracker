package symbols

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Marv963/CryptoTracker/app/pkg/websocket"
)

type NotifyHandler struct {
	Server *websocket.Server
}

func NewNotifyHandler(server *websocket.Server) *NotifyHandler {
	return &NotifyHandler{
		Server: server,
	}
}

func (n *NotifyHandler) NotifyPriceChange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Hier wird einfach der Body als []byte gelesen
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	// könnten Sie es hier schnell validieren, ohne es zu deserialisieren.
	var js map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &js); err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	// Pair aus dem validierten JSON extrahieren für den Broadcast.
	pair, ok := js["pair"].(string)
	if !ok {
		http.Error(w, "invalid or missing pair", http.StatusBadRequest)
		return
	}

	n.Server.Broadcast(pair, bodyBytes)
}
