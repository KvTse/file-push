package config

type RedisConfig struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBIndex  int    `mapstructure:"dbIndex" json:"dbIndex" yaml:"dbIndex"`
}
