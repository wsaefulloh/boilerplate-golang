package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakejorge/boilerplate-golang/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response.OK(c, "Service is healthy and running", gin.H{
		"status":    "UP",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
