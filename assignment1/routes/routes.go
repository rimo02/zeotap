package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rimo02/zeotap/assignment1/controllers"
)

var RegisterRoutes = func(c *gin.Engine) {
	c.POST("/create", controllers.CreateRule)
	c.POST("/combine", controllers.CombineRules)
	c.GET("/evaluate", controllers.EvaluateRule)
}