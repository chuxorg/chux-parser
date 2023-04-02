package config

import (
	"fmt"

	"github.com/spf13/viper"
)


type ParserConfig struct {
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
	AWS struct {
		BucketName   string `mapstructure:"bucketName"`
		DownloadPath string `mapstructure:"downloadPath"`
	} `mapstructure:"aws"`
	Auth struct {
		IssuerURL string `mapstructure:"issuerUrl"`
		TokenURL  string `mapstructure:"tokenUrl"`
	} `mapstructure:"auth"`
	BizConfig struct {
		DataStores []struct {
			DataStore struct {
				Mongo struct {
					Target         string        `mapstructure:"target"`
					URI            string        `mapstructure:"uri"`
					Timeout        int 	         `mapstructure:"timeout"`
					DatabaseName   string        `mapstructure:"databaseName"`
					CollectionName string        `mapstructure:"collectionName"`
				} `mapstructure:"mongo"`
			} `mapstructure:"dataStore"`
		} `mapstructure:"dataStores"`
	} `mapstructure:"bizConfig"`
}



func LoadConfig(env string) (*ParserConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fmt.Sprintf("config.%s", env)) // e.g., config.development or config.production
	viper.AddConfigPath(".")                           // Look for config files in the current directory
	viper.AddConfigPath("./config")                    // Look for config files in the config directory
	viper.AddConfigPath("../../config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	var cfg ParserConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	// Make sure the DataStores and DataStoreMap are initialized
	if len(cfg.BizConfig.DataStores) == 0 {
		cfg.BizConfig.DataStores = make([]struct {
			DataStore struct {
				Mongo struct {
					Target         string        `mapstructure:"target"`
					URI            string        `mapstructure:"uri"`
					Timeout        int            `mapstructure:"timeout"`
					DatabaseName   string        `mapstructure:"databaseName"`
					CollectionName string        `mapstructure:"collectionName"`
				} `mapstructure:"mongo"`
			} `mapstructure:"dataStore"`
		}, 0)
	}

	return &cfg, nil
}
