package appcontext

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/Marv963/CryptoTracker/app/pkg/config"
	"github.com/Marv963/CryptoTracker/app/pkg/database/influxdb"
)

type AppContext struct {
	InfluxClient *influxdb.InfluxDBClient
	Redis        *redis.Client
	Logger       *log.Logger
}

func NewAppContext(cfg *config.Config) (*AppContext, error) {
	influxClient, err := influxdb.NewClient(*cfg.InfluxDB)
	if err != nil {
		return nil, fmt.Errorf("error cannot create InfluxDB Client: %w", err)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	logger := log.New(os.Stdout, "APP_LOG: ", log.LstdFlags)

	return &AppContext{
		Logger:       logger,
		InfluxClient: influxClient,
		Redis:        redis,
	}, nil
}

func (a AppContext) Close() error {
	a.InfluxClient.Close()
	err := a.Redis.Close()
	if err != nil {
		return fmt.Errorf("error cannot close redis: %w", err)
	}
	return nil
}
