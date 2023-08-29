package config

type Server struct {
	RedisConfig RedisConfig `mapstructure:"redis" json:"redis" yaml:"redis"`
	KafkaConfig KafkaConfig `mapstructure:"kafka" json:"kafka" yaml:"kafka"`
}
