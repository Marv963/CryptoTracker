package config

type RabbitMQ struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	URL      string
	Exchange string `mapstructure:"exchange"`
}
type InfluxDB struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Token    string `mapstructure:"token"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	URL      string
	Org      string `mapstructure:"org"`
	Bucket   string `mapstructure:"bucket"`
}

type Redis struct {
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	URL      string
}

type Websocket struct {
	URL string `mapstructure:"url"`
}

type Frontend struct {
	URL string `mapstructure:"url"`
}

type Config struct {
	RabbitMQ  *RabbitMQ
	InfluxDB  *InfluxDB
	Redis     *Redis
	Websocket *Websocket
	Frontend  *Frontend
}
