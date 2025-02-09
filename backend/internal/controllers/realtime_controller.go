package controllers

import (
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type RealTimeController struct {
	mu      sync.Mutex
	clients map[chan string]bool
}

func NewRealTimeController() *RealTimeController {
	return &RealTimeController{
		clients: make(map[chan string]bool),
	}
}

// RegisterRoutes can be called in main.go to add SSE endpoints
func (rtc *RealTimeController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/updates", rtc.HandleSSE)
	rg.POST("/broadcast", rtc.BroadcastMessage)
}

// HandleSSE => GET /api/updates
func (rtc *RealTimeController) HandleSSE(c *gin.Context) {
	messageChan := make(chan string)

	// Add this client
	rtc.mu.Lock()
	rtc.clients[messageChan] = true
	rtc.mu.Unlock()

	// Remove this client when done
	defer func() {
		rtc.mu.Lock()
		delete(rtc.clients, messageChan)
		rtc.mu.Unlock()
		close(messageChan)
	}()

	c.Stream(func(w io.Writer) bool {
		// Wait for a message
		msg, ok := <-messageChan
		if !ok {
			return false
		}
		c.SSEvent("update", msg)
		return true
	})
}

// BroadcastMessage => POST /api/broadcast
// Send the same "message" to all connected SSE clients.
func (rtc *RealTimeController) BroadcastMessage(c *gin.Context) {
	var payload struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rtc.mu.Lock()
	defer rtc.mu.Unlock()

	for ch := range rtc.clients {
		ch <- payload.Message
	}

	c.JSON(http.StatusOK, gin.H{"status": "broadcasted"})
}
