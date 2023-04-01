package version

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	Ver       string
	GitHash   string
	BuildDate string
)

func Version(c *gin.Context) {
	res := map[string]interface{}{
		"version":   Ver,
		"gitHash":   GitHash,
		"buildDate": BuildDate,
	}
	log.Println("Version called..")
	c.JSON(http.StatusOK, res)
}
