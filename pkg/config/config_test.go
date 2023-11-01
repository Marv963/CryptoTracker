package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	// Save current environment variable values
	originalRabbitMQHost := os.Getenv("RABBITMQ_HOST")
	originalRabbitMQPort := os.Getenv("RABBITMQ_PORT")

	// Set environment variables for testing
	os.Setenv("RABBITMQ_HOST", "rabbitmq")
	os.Setenv("RABBITMQ_PORT", "9090")

	defer func() {
		// Restore the original environment variable values after the test
		os.Setenv("RABBITMQ_HOST", originalRabbitMQHost)
		os.Setenv("RABBITMQ_PORT", originalRabbitMQPort)
	}()

	// Create a temporary configuration file
	configContent := `
    rabbitmq:
      host: 127.0.0.1
      port: 5672
    influxdb:
      host: 127.0.0.1
      port: 8086
    redis:
      host: 127.0.0.1
      port: 6379
    `
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(configContent))
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Call the Read function
	config, err := Read(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	viper.Reset()
	// Assertions
	assert := assert.New(t)
	assert.Equal("rabbitmq:9090", config.RabbitMQ.URL, "RabbitMQ URL should be correct")
	assert.Equal("http://127.0.0.1:8086", config.InfluxDB.URL, "InfluxDB URL should be correct")
	assert.Equal("127.0.0.1:6379", config.Redis.URL, "Redis URL should be correct")
}
