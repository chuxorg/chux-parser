package config

import (
	"fmt"

	"github.com/chuxorg/chux-models/config"
	mcfg "github.com/chuxorg/chux-models/config"
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
	
	BizObjConfig config.BizObjConfig `mapstructure:"dataStores"`
}



func LoadConfig(env string) (*ParserConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fmt.Sprintf("config.%s.yaml", env)) // e.g., config.development or config.production
	viper.AddConfigPath("../config") // Look for config file in the parent directory/config
	
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	var cfg ParserConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	// Initialize BizObjConfig.DataStores.DataStoreMap if it's not already set
	if cfg.BizObjConfig.DataStores.DataStoreMap == nil {
		cfg.BizObjConfig.DataStores.DataStoreMap = make(map[string]mcfg.DataStoreConfig)
	}

	
	return &cfg, nil
}

func GetBizObjConfig(cfg ParserConfig) mcfg.BizObjConfig {
	bizObjConfig := mcfg.BizObjConfig{
		Logging: cfg.Logging,
		DataStores: cfg.BizObjConfig.DataStores,
	}
	return bizObjConfig
}
