package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/reww406/linetracker/config"
	"github.com/reww406/linetracker/internal/metro"
	"github.com/reww406/linetracker/internal/station"
	"github.com/reww406/linetracker/internal/store"
	"github.com/reww406/linetracker/internal/train"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *gin.Engine
}

var ddbClient *dynamodb.Client

func (s *Server) getNextTrains(c *gin.Context) {
	// Get query parameters
	lineCode := c.Query("line_code")
	locationCode := c.Query("location_code")
	direction := c.Query("direction")

	req := train.GetNextTrainsRequest{
		LineCode:     metro.LineCode(lineCode),
		LocationCode: locationCode,
		Direction:    direction,
	}

	result, err := train.GetTrainPredictions(c, ddbClient, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get trains",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trains": result,
	})
}

func (s *Server) getStations(c *gin.Context) {
	stationList, err := station.ListStations(c, ddbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get stations",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"stations": stationList,
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes group
	v1 := s.router.Group("/api/v1")
	{
		// Routes
		// api/v1/stations
		v1.GET("/stations", func(c *gin.Context) {
			s.getStations(c)
		})

		// api/v1/trains?line_code=RD&line_code=BL&location_code=A01&direction_code=B01
		v1.GET("/trains", func(c *gin.Context) {
			s.getNextTrains(c)
		})
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

	client, err := store.InitDB()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to connect to DDB.")
	}
	ddbClient = client
	go train.PollTrainPredictions(ddbClient)
	server := CreateGinServer()
	if err := server.router.Run(fmt.Sprintf(":%d", config.BindingPort)); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to start server.")
	}
}
