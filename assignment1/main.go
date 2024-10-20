package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rimo02/zeotap/assignment1/database"
	"github.com/rimo02/zeotap/assignment1/routes"
	"net/http"
)

func init() {
	godotenv.Load()
	database.InitializeConnections()
}

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("./templates/index.html")

	routes.RegisterRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	err := r.Run(":8000")
	if err != nil {
		fmt.Printf("Error running at port 8000\n")
	}
}
