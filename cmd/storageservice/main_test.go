package main

import (
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInfluxClient is a mock for InfluxDBClient
type MockInfluxClient struct {
	mock.Mock
}

func (m *MockInfluxClient) WriteData(measurement string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) error {
	args := m.Called(measurement, tags, fields, timestamp)
	return args.Error(0)
}

func (m *MockInfluxClient) Close() {}

// Test the processMessage function
func TestProcessMessage(t *testing.T) {
	// Mock the InfluxDB client
	influxMock := new(MockInfluxClient)

	app := &App{
		InfluxClient: influxMock,
		RabbitMQ:     nil, // This is not used in processMessage, so it's okay to be nil
	}

	// Example delivery message
	d := amqp091.Delivery{
		Body: []byte(`{
			"timestamp": 1675234567890,
			"symbol": "BTC/USD",
   		"high": 45000.23,
      "low": 43000.12,
      "bid": 44500.5,
      "bidVolume": 20.5,
      "ask": 44950.25,
      "askVolume": 15.3,
      "vwap": 44000.18,
      "open": 42000.0,
      "close": 44800.56,
      "last": 44850.1,
      "previousClose": 42000.0,
      "change": 2800.56,
      "percentage": 6.67,
      "average": 43500.25,
      "baseVolume": 1000.2,
      "quoteVolume": 43500250,
      "info": {
        "open": "42000.0",
        "high": "45000.23",
        "low": "41500.12",
        "last": "44850.1",
        "volume": "1000.2",
        "vwap": "43500.25",
        "bid": "44500.5",
        "ask": "44950.25",
        "side": "buy",
        "open_24": "40000.0",
        "percent_change_24": "12.125"
      }
		}`),
	}

	// Set up expectations
	influxMock.On("WriteData", "crypto_data", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := app.processMessage("kraken", d)

	assert.Nil(t, err)
	influxMock.AssertExpectations(t)
}

func TestProcessMessage_MissingTimestamp(t *testing.T) {
	// Mock the InfluxDB client
	influxMock := new(MockInfluxClient)

	app := &App{
		InfluxClient: influxMock,
		RabbitMQ:     nil, // This is not used in processMessage, so it's okay to be nil
	}

	// Example delivery message without timestamp
	d := amqp091.Delivery{
		Body: []byte(`{
			"symbol": "BTC/USD"
		}`),
	}

	err := app.processMessage("kraken", d)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing Timestamp") // This checks if the error message contains the string "missing Timestamp"
}

func TestProcessMessage_TimestampZero(t *testing.T) {
	// Mock the InfluxDB client
	influxMock := new(MockInfluxClient)

	app := &App{
		InfluxClient: influxMock,
		RabbitMQ:     nil,
	}

	// Example delivery message with timestamp being 0
	d := amqp091.Delivery{
		Body: []byte(`{
			"timestamp": 0,
			"symbol": "BTC/USD"
		}`),
	}

	err := app.processMessage("kraken", d)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing Timestamp")
}
