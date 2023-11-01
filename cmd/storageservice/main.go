package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Marv963/CryptoTracker/app/pkg/config"
	influxdb "github.com/Marv963/CryptoTracker/app/pkg/database/influxdb"
	"github.com/Marv963/CryptoTracker/app/pkg/messaging"
	model "github.com/Marv963/CryptoTracker/app/pkg/model/ccxt"
	"github.com/rabbitmq/amqp091-go"
)

type InfluxDBClientInterface interface {
	WriteData(measurement string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) error
	Close()
}

type App struct {
	InfluxClient InfluxDBClientInterface
	RabbitMQ     *messaging.RabbitMQ
}

func NewApp(cfg *config.Config) (*App, error) {
	influxClient, err := influxdb.NewClient(*cfg.InfluxDB)
	if err != nil {
		log.Fatalf("Failed to create InfluxDB client: %s", err)
	}
	defer influxClient.Close()

	rmqClient, err := messaging.NewRabbitMQ(*cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	cfg = nil
	return &App{InfluxClient: influxClient, RabbitMQ: rmqClient}, nil
}

func (a *App) shutdown() {
	a.InfluxClient.Close()
	a.RabbitMQ.Channel.Close()
}

func (app *App) processMessage(exchange string, d amqp091.Delivery) error {
	// Process the message...
	var msg model.Message
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		log.Printf("Failed to parse message data: %s : message: %s", err, d.Body)
		return err
	}

	if msg.Timestamp == 0 {
		return fmt.Errorf("missing Timestamp")
	}

	// Convert the timestamp from the message
	timestamp := time.Unix(0, msg.Timestamp*int64(time.Millisecond))

	// Tags and fields for InfluxDB
	tags := map[string]string{
		"exchange": exchange, // Replace this with the actual exchange name
		"pair":     msg.Symbol,
	}
	fields := map[string]interface{}{
		"high":                   msg.High,
		"low":                    msg.Low,
		"bid":                    msg.Bid,
		"bidVolume":              msg.BidVolume,
		"ask":                    msg.Ask,
		"askVolume":              msg.AskVolume,
		"vwap":                   msg.Vwap,
		"open":                   msg.Open,
		"close":                  msg.Close,
		"last":                   msg.Last,
		"previousClose":          msg.PreviousClose,
		"change":                 msg.Change,
		"percentage":             msg.Percentage,
		"average":                msg.Average,
		"base_volume":            msg.BaseVolume,
		"quote_volume":           msg.QuoteVolume,
		"info_open":              msg.Info.Open,
		"info_high":              msg.Info.High,
		"info_low":               msg.Info.Low,
		"info_last":              msg.Info.Last,
		"info_volume":            msg.Info.Volume,
		"info_vwap":              msg.Info.Vwap,
		"info_bid":               msg.Info.Bid,
		"info_ask":               msg.Info.Ask,
		"info_side":              msg.Info.Side,
		"info_open_24":           msg.Info.Open24,
		"info_percent_change_24": msg.Info.PercentChange24,
	}

	// Write data to InfluxDB
	err = app.InfluxClient.WriteData("crypto_data", tags, fields, timestamp)
	if err != nil {
		log.Printf("Failed to write data to InfluxDB: %s", err)
		return err
	}

	return nil
}

func main() {
	cfg, err := config.Read("")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to read configuration", err)
	}
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %s", err)
	}
	defer app.shutdown()

	exchanges := []string{"kraken", "bitstamp"}

	// Create different Go routines for each exchange
	for _, exchange := range exchanges {
		queue := fmt.Sprintf("queue_%s_prices", exchange)
		// Here, each key could represent an exchange
		go func(exchange string) {
			// Queue binding with different routing keys
			routingKey := fmt.Sprintf("%s.prices", exchange)
			err := app.RabbitMQ.BindQueue(app.RabbitMQ.Exchange, queue, routingKey)
			if err != nil {
				log.Fatalf("Failed to bind queue: %s", err)
			}

			// Receive and process messages
			deliveries, err := app.RabbitMQ.Receive(queue)
			if err != nil {
				log.Fatalf("Failed to start consuming messages: %s", err)
			}

			for d := range deliveries {
				// Process message...
				err := app.processMessage(exchange, d)
				if err != nil {
					log.Printf("Failed to process message: %s", err)
				}
				d.Ack(false)
			}
		}(exchange)
	}

	forever := make(chan struct{})
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
