package main

import (
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"

	"example.com/se/config"
	"example.com/se/controller/candidates"
	"example.com/se/controller/elections"
	"example.com/se/controller/genders"
	"example.com/se/controller/users"
	"example.com/se/controller/votes"
	"example.com/se/entity"
	"example.com/se/metrics"
	"example.com/se/middlewares"
)

const (
	PORT           = "8000"
	PushGatewayURL = "http://pushgateway:9091"
	JobName        = "online_voting_system"
)

func main() {
	// DB
	config.ConnectionDB()
	config.SetupDatabase()

	// Register metrics
	metrics.RegisterMetrics()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// HTTP metrics middleware
	r.Use(func(c *gin.Context) {
		timer := prometheus.NewTimer(metrics.HTTPRequestDurationSeconds.
			WithLabelValues(c.Request.Method, c.FullPath()))
		defer timer.ObserveDuration()

		c.Next()

		metrics.HTTPRequestsTotal.
			WithLabelValues(c.Request.Method, c.FullPath(), http.StatusText(c.Writer.Status())).
			Inc()
	})

	// Prometheus scrape endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routes
	api := r.Group("/api")
	{
		api.POST("/signup", users.SignUp)
		api.POST("/signin", users.SignIn)
		api.GET("/genders", genders.GetAll)

		api.Use(middlewares.Authorizes())
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

	// Optional: periodically push snapshot metrics
	go func() {
		for {
			pushSnapshotMetrics()
			time.Sleep(30 * time.Second)
		}
	}()

	r.Run(":" + PORT)
}

func pushSnapshotMetrics() {
	db := config.DB()
	var userCount, voteCount int64
	db.Model(&entity.Users{}).Count(&userCount)
	db.Model(&entity.Votes{}).Count(&voteCount)

	// Update gauges
	metrics.UsersTotal.Set(float64(userCount))

	// We defined VotesCreatedTotal as Counter; for snapshot use Inc
	for i := int64(0); i < voteCount; i++ {
		metrics.VotesCreatedTotal.Inc()
	}

	// Push to Pushgateway
	_ = push.New(PushGatewayURL, JobName).
		Collector(metrics.UsersTotal).
		Collector(metrics.VotesCreatedTotal).
		Push()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
