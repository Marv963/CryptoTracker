package influxdb

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/Marv963/CryptoTracker/app/pkg/config"
)

type InfluxDBClient struct {
	Client influxdb2.Client
	Bucket string
	org    string
}

// NewClient erstellt einen neuen InfluxDB-Client.
func NewClient(cfg config.InfluxDB) (*InfluxDBClient, error) {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)
	return &InfluxDBClient{Client: client, Bucket: cfg.Bucket, org: cfg.Org}, nil
}

// WriteData schreibt Daten in die InfluxDB.
func (c *InfluxDBClient) WriteData(measurement string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) error {
	writeAPI := c.Client.WriteAPIBlocking(c.org, c.Bucket)
	p := influxdb2.NewPoint(measurement, tags, fields, timestamp)
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return fmt.Errorf("error: Failed to write data to InfluxDB: %w", err)
	}
	return nil
}

// Close beendet die Client-Verbindung.
func (c *InfluxDBClient) Close() {
	c.Client.Close()
}

func (c *InfluxDBClient) replaceBucketPlaceholder(query string) string {
	return strings.ReplaceAll(query, "{bucket}", c.Bucket)
}

// QueryData fragt Daten von der InfluxDB ab.
func (c *InfluxDBClient) queryData(query string) (*api.QueryTableResult, error) {
	queryAPI := c.Client.QueryAPI(c.org)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error: Failed to query data from InfluxDB : %v", err)
	}
	return result, nil
}

// QueryWrapper
func (c *InfluxDBClient) Query(queryPath string, args ...interface{}) (result *api.QueryTableResult, err error) {
	// Lese die Query aus der Datei
	queryBytes, err := os.ReadFile(queryPath)
	if err != nil {
		return nil, fmt.Errorf("error: Reading query file at path [%s] failed: %w", queryPath, err)
	}

	// Ersetze weitere Platzhalter
	queryString := fmt.Sprintf(string(queryBytes), args...)

	// Ersetze den Bucket-Namen und eventuelle weitere Platzhalter
	queryString = c.replaceBucketPlaceholder(queryString)

	// Führe die Query aus
	result, err = c.queryData(queryString)
	if err != nil {
		return nil, fmt.Errorf("error: Querying InfluxDB at [%s] failed: %w", queryPath, err)
	}

	return result, nil
}
