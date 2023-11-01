package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	symbols "github.com/Marv963/CryptoTracker/app/cmd/websocket/handler"
	"github.com/Marv963/CryptoTracker/app/pkg/websocket"
)

func main() {
	server := websocket.NewServer()
	server.Start()
	defer server.Shutdown()

	http.HandleFunc("/ws", server.HandleWS)

	symbolController := symbols.NewNotifyHandler(server)
	http.HandleFunc("/notify", symbolController.NotifyPriceChange)
	http.ListenAndServe(":3008", nil)

	// Setup signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Attempt a graceful shutdown.
	server.Shutdown()
	fmt.Println("Shutting down")
	os.Exit(0)
}
