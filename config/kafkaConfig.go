package config

type KafkaConfig struct {
	Brokers         string `mapstructure:"brokers" json:"brokers" yaml:"brokers"`
	Topic           string `mapstructure:"topic" json:"topic" yaml:"topic"`
	ConsumerGroupId string `mapstructure:"consumer-group-id" json:"consumerGroupId" yaml:"consumer-group-id"`
}
