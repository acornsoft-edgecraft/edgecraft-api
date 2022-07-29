// Package config - Configuration
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db/postgresdb"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// Config - Represents the configuration
type Config struct {
	API *api.Config        `yaml:"api"`
	DB  *postgresdb.Config `yaml:"database"`
	// Mail       *mail.Config       `yaml:"mail"`
	// K8sGateway *k8sclient.Config  `yaml:"k8sgateway"`
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// Load - Load the configuration from file
func Load() (*Config, error) {

	// Search config files in config directory with name "config.yaml" (without extension).
	viper.AddConfigPath("./conf")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// yaml의 속성 ex) database.database_name 등을 update 위한 환경변수 설정시 . 대신 '_' 사용할 수 있게 한다.
	// ex)  database.database_name 속성은  DATABASE_DATABASE_NAME 환경변수 설정하면 값이 override 됨
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// how to bind the env or the flag
	// viper.BindPFlag("port", serverCmd.Flags().Lookup("port")) // flag-viper binding
	// viper.BindEnv("home") // binding with env HOME

	var conf = Config{}
	var err error

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			logger.Errorf("Could not load configuration file: %s", err.Error())
			os.Exit(0)
		} else {
			// Config file was found but another error was produced
		}
	}

	logger.Info("Using config file:", viper.ConfigFileUsed())

	// Unmarshal to instance
	// viper.Unmarshal(&conf)
	viper.Unmarshal(&conf, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "yaml"
	})

	fmt.Printf("config info : %v", conf)

	return &conf, err
}
