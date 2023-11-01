package websocket

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

func TestServerHandleWS(t *testing.T) {
	s := NewServer()
	s.Start()
	defer s.Shutdown()

	// Create a HTTP test server
	server := httptest.NewServer(http.HandlerFunc(s.HandleWS))
	defer server.Close()

	// Create a websocket client connection
	conn, _, err := gorillaws.DefaultDialer.Dial("ws"+server.URL[4:], nil)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	// Delay to ensure that the connection has been established
	time.Sleep(100 * time.Millisecond)

	// Check if a connection has been added
	if count := atomic.LoadInt64(&s.activeConns); count != 1 {
		t.Fatalf("Expected 1 active connection, got %d", count)
	}

	// TODO: Add more tests here, e.g. for sending/receiving messages, managing subscriptions, etc.

	// Delay to ensure that the connection has been removed
	time.Sleep(100 * time.Millisecond)

	conn.Close()
	time.Sleep(100 * time.Millisecond)

	if count := atomic.LoadInt64(&s.activeConns); count != 0 {
		t.Fatalf("Expected 0 active connections after disconnected, got %d", count)
	}
}
