package symbols

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Marv963/CryptoTracker/app/pkg/websocket"
)

func TestNotifyPriceChange(t *testing.T) {
	// Create a WebSocket server
	wsserver := websocket.NewServer()
	defer wsserver.Shutdown()

	// HTTP Handler for broadcasting messages
	handler := NewNotifyHandler(wsserver)

	req, err := http.NewRequest("POST", "/notify", strings.NewReader(`{"pair":"BTC/USD","price":"50000"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.NotifyPriceChange)

	handlerFunc.ServeHTTP(rr, req)

	// Check the return status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
