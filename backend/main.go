package main

import (
	"example.com/se/config"
	"example.com/se/controller/candidates"
	"example.com/se/controller/elections"
	"example.com/se/controller/votes"
	"example.com/se/controller/genders"
	"example.com/se/controller/users"
	"example.com/se/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)


const (
    PORT           = "8000"
    PushGatewayURL = "http://pushgateway:9091"
    JobName        = "online_voting_system"
)

func main() {

	// open connection database

	config.ConnectionDB()

	// Generate databases

	config.SetupDatabase()

	r := gin.Default()

	r.Use(CORSMiddleware())

	// Auth Route

	api := r.Group("/api")
	{
		api.POST("/signup", users.SignUp)
		api.POST("/signin", users.SignIn)
		api.GET("/genders", genders.GetAll)
		// Protected Routes
		api.Use(middlewares.Authorizes()) // ✅ ใส่ middleware ที่นี่
		{
			api.PUT("/user/:id", users.Update)
			api.GET("/users", users.GetAll)
			api.GET("/user/:id", users.Get)
			api.DELETE("/user/:id", users.Delete)

			api.POST("/candidate", candidates.Create)
			api.PUT("/candidate/:id", candidates.Update)
			api.GET("/candidates", candidates.GetAll)
			api.GET("/candidate/:id", candidates.Get)
			api.DELETE("/candidate/:id", candidates.Delete)

			api.POST("/election", elections.Create)
			api.PUT("/election/:id", elections.Update)
			api.GET("/elections", elections.GetAll)
			api.GET("/election/:id", elections.Get)
			api.DELETE("/election/:id", elections.Delete)

			api.GET("/votes", votes.GetAll)
			api.GET("/vote/:id", votes.Get)
			api.POST("/vote", votes.CreateVote)
		}
		
	}

	

	r.GET("/", func(c *gin.Context) {

		c.String(http.StatusOK, "API RUNNING... PORT: %s", PORT)

	})
	
	// Run the server

	r.Run(":" + PORT)


}

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {

			c.AbortWithStatus(204)

			return

		}

		c.Next()

	}

}
