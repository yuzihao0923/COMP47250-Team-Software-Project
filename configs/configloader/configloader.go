package configloader

type Config struct {
	Brokers []BrokerConfig
}

type BrokerConfig struct {
	ID      string
	Address string
}

type ConfigLoader interface {
	LoadConfig() (*Config, error)
}
