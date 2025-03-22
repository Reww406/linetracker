package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/reww406/linetracker/config"	
	"github.com/reww406/linetracker/internal/store"	
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *gin.Engine
}

//func (s *Server) getStations(c *gin.Context) {
//    c.JSON(http.StatusOK, gin.H{
//        "stations": []string{"station1", "station2"},
//    })
//}
//
//func (s *Server) Run(addr string) error {
//    return s.router.Run(addr)
//}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes group
	v1 := s.router.Group("/api/v1")
	{
		// Routes
		v1.GET("")
	}
}

func CreateGinServer() *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	server := &Server{
		router: router,
	}

	server.setupRoutes()
	return server
}

func main() {
	log := config.GetLogger()
	config := config.LoadConfig()

  _, err := store.InitDB()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to open DB.")
	}

	server := CreateGinServer()
	if err := server.router.Run(fmt.Sprintf(":%d", config.BindingPort)); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to start server.")
	}
}
