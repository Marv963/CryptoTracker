package main

import (
	"log"
	"os"
	"time"

	"github.com/Marv963/CryptoTracker/app/cmd/datarefresher/jobs"
	"github.com/Marv963/CryptoTracker/app/pkg/appcontext"
	"github.com/Marv963/CryptoTracker/app/pkg/config"
)

func main() {
	cfg, err := config.Read("")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to read configuration", err)
		os.Exit(2)
	}

	appContext, err := appcontext.NewAppContext(cfg)
	if err != nil {
		appContext.Logger.Fatalf("Failed to initialize app context: %s", err)
		os.Exit(2)
	}
	defer appContext.Close()

	// Instanziieren Sie Ihre Jobs
	job := jobs.NewJobController(appContext, cfg.Websocket.URL)
	job.Start()
	defer job.Shutdown()
	go job.FetchDataWithInterval(job.FetchSymbolQuotes, 1*time.Minute)
	go job.FetchDataWithInterval(job.FetchSymbolPriceChanges, 5*time.Second)

	// Hauptgoroutine am Leben halten, um die Hintergrundgoroutine weiter auszuführen
	forever := make(chan struct{})
	log.Printf(" [*] Executing Jobs. To exit press CTRL+C")
	<-forever
}
