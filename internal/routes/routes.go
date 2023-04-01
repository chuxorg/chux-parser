package routes

import (
	"log"
	"net/http"

	"chux-parser/internal/auth"
	"chux-parser/internal/groups"
	"idc-okta-api/internal/version"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Server() *gin.Engine {

	router := gin.Default()
	router.Use(cors.Default())

	//--Route for getting all groups
	grpSvc := groups.NewGroupService()

	v1 := router.Group("/api/v1")
	{
		groups := v1.Group("/groups", auth.AuthMiddleware())
		{
			groups.GET("/", grpSvc.Fetch)
			groups.GET(":groupId/users", grpSvc.GetGroupUsers)
		}
		helpers := v1.Group("/")
		{
			helpers.GET("ping", Ping)
			helpers.GET("version", version.Version)
			// helpers.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

	}

	log.Println("These are the configured routes.")
	log.Println(router.Routes())
	//--return router
	return router
}

// Ping godoc
// @Summary Shows the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/ping [get]
func Ping(c *gin.Context) {
	res := map[string]interface{}{
		"data": "idc-okta-api Service is up and running",
	}
	log.Println("Health check called..")
	c.JSON(http.StatusOK, res)
}
