package config

import (
	"os"
	"testing"

	"github.com/chuxorg/chux-models/models/products"
	"github.com/chuxorg/chux-parser/config"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/assert"
)

// TestNew tests the New function with different options.
func TestNew(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	_cfg, err := config.LoadConfig("test")
	assert.Nil(t, _cfg)
	assert.NotNil(t, _cfg.Logging)
	assert.NotNil(t, _cfg.DataStores)
	assert.Equal(1, len(_cfg.DataStores.DataStoreMap))
	
	product := products.WithBizObjConfig(_cfg.BizConfig)
	assert.NotNil(t, product)
	assert.Equal(t, "chux-cprs", product.GetDatabaseName())
	assert.Equal(t, "products", product.GetCollectionName())
	assert.Equal(t, "mongodb://localhost:27017", product.GetURI())
}