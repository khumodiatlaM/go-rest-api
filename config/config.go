package config

import (
	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Broker string `mapstructure:"KAFKA_BROKER"`
	Topic  string `mapstructure:"KAFKA_TOPIC"`
}
type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	APIPort    string `mapstructure:"API_PORT"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	Kafka      KafkaConfig
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (Config, error) {
	var config Config
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	err = viper.Unmarshal(&config.Kafka)
	return config, err
}
