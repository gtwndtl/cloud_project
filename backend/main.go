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

const PORT = "8000"

func main() {

	// open connection database

	config.ConnectionDB()

	// Generate databases

	config.SetupDatabase()

	r := gin.Default()

	r.Use(CORSMiddleware())

	// Auth Route

	r.POST("/signup", users.SignUp)

	r.POST("/signin", users.SignIn)

	router := r.Group("/")

	{

		router.Use(middlewares.Authorizes())

		// User Route

		router.PUT("/user/:id", users.Update)

		router.GET("/users", users.GetAll)

		router.GET("/user/:id", users.Get)

		router.DELETE("/user/:id", users.Delete)

		router.POST("/candidate", candidates.Create)

		router.PUT("/candidate/:id", candidates.Update)

		router.GET("/candidates", candidates.GetAll)

		router.GET("/candidate/:id", candidates.Get)

		router.DELETE("/candidate/:id", candidates.Delete)

		router.POST("/election", elections.Create)

		router.PUT("/election/:id", elections.Update)

		router.GET("/elections", elections.GetAll)

		router.GET("/election/:id", elections.Get)

		router.DELETE("/election/:id", elections.Delete)

		router.GET("/votes", votes.GetAll)

		router.GET("/vote/:id", votes.Get)

		router.POST("/vote", votes.CreateVote)

	}

	r.GET("/genders", genders.GetAll)

	r.GET("/", func(c *gin.Context) {

		c.String(http.StatusOK, "API RUNNING... PORT: %s", PORT)

	})

	// Run the server

	r.Run("localhost:" + PORT)

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
