package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	r := gin.New()

	// Middleware for logging
	r.Use(func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// After request
		latency := time.Since(startTime)
		status := c.Writer.Status()

		// Log request details
		logrus.WithFields(logrus.Fields{
			"client_ip":  c.ClientIP(),
			"timestamp":  startTime.Format(time.RFC3339),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     status,
			"latency":    latency,
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP request")
	})

	r.GET("/", func(c *gin.Context) {
		// Get system time
		currentTime := time.Now().Format("2006-01-02 15:04:05")

		// Get hostname
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}

		// Get IP address
		ipAddr := "unknown"
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						ipAddr = ipNet.IP.String()
						break
					}
				}
			}
		}

		// Get client IP
		clientIP := c.ClientIP()

		// Get custom environment variable
		imageVersion := os.Getenv("imgVer")
		if imageVersion == "" {
			imageVersion = "not set"
		}

		// Get binary name and extract the last part
		binaryName := os.Args[0] // 获取完整路径
		binaryFileName := path.Base(binaryName)

		// Construct response
		response := fmt.Sprintf(
			"Go APP Demo\n"+
				"=================\n"+
				"SystemTime : %s\n"+
				"HostName   : %s\n"+
				"IPAddress  : %s\n"+
				"ClientIP   : %s\n"+
				"ImageVersion: %s\n"+
				"AppName : %s\n",
			currentTime, hostName, ipAddr, clientIP, imageVersion, binaryFileName)

		c.String(http.StatusOK, response)
	})

	// Add HEAD request support for the root route
	r.HEAD("/", func(c *gin.Context) {
		c.Status(http.StatusOK) // Respond with 200 OK for HEAD requests
	})

	r.GET("/health", func(c *gin.Context) {
		healthStatus := "Healthy"
		c.JSON(http.StatusOK, gin.H{
			"status": healthStatus,
		})
	})

	// Add HEAD request support for the health route
	r.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK) // Respond with 200 OK for HEAD requests
	})

	// Log service startup
	logrus.Info("Starting server on :8080")

	if err := r.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("Couldn't listen")
	}
}
