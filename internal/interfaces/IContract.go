package interfaces

import "github.com/gin-gonic/gin"

type IContract interface {
	Fetch(c *gin.Context)
	Ping(c *gin.Context)
}
