package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var config *Config

type Config struct {
	RequestsPerIP    int           `mapstructure:"REQUESTS_PER_IP"`
	RequestsPerToken int           `mapstructure:"REQUESTS_PER_TOKEN"`
	BlockDuration    time.Duration `mapstructure:"BLOCK_DURATION"`
}

func NewConfig() *Config {
	return config
}

func init() {
	var err error
	config, err = loadConfig()
	if err != nil {
		panic(fmt.Sprintf("Erro ao carregar configurações: %v", err))
	}
}

func loadConfig() (*Config, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Habilita o carregamento de variáveis de ambiente
	viper.AutomaticEnv()

	// Tenta ler o arquivo .env
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente")
	}

	// Define valores padrão
	viper.SetDefault("REQUESTS_PER_IP", 10)
	viper.SetDefault("REQUESTS_PER_TOKEN", 100)
	viper.SetDefault("BLOCK_DURATION", "5m")

	// Cria uma nova instância de Config
	config := &Config{}

	// Carrega as variáveis de ambiente
	requestsPerIP := viper.GetInt("REQUESTS_PER_IP")
	requestsPerToken := viper.GetInt("REQUESTS_PER_TOKEN")
	blockDurationStr := viper.GetString("BLOCK_DURATION")

	blockDuration, err := time.ParseDuration(blockDurationStr)

	if err != nil {
		return nil, fmt.Errorf("erro ao analisar BLOCK_DURATION: %v", err)
	}

	config.RequestsPerIP = requestsPerIP
	config.RequestsPerToken = requestsPerToken
	config.BlockDuration = blockDuration

	// Validação das configurações obrigatórias
	if config.RequestsPerIP <= 0 {
		return nil, fmt.Errorf("REQUESTS_PER_IP deve ser maior que zero")
	}

	if config.RequestsPerToken <= 0 {
		return nil, fmt.Errorf("REQUESTS_PER_TOKEN deve ser maior que zero")
	}

	if config.BlockDuration <= 0 {
		return nil, fmt.Errorf("BLOCK_DURATION deve ser maior que zero")
	}
	return config, nil
}
