package controllers

import (
	"io"       // For writing Server-Sent Event responses.
	"net/http" // Provides HTTP status codes.
	"sync"     // Used for mutex locks when accessing shared data.

	"github.com/gin-gonic/gin" // Gin framework for HTTP routing.
)

// RealTimeController manages Server-Sent Events (SSE) for real-time updates.
// It maintains a set of client channels to which messages can be broadcast.
type RealTimeController struct {
	mu      sync.Mutex           // Mutex to ensure safe concurrent access to the clients map.
	clients map[chan string]bool // Map storing all active client channels.
}

// NewRealTimeController creates and initializes a new RealTimeController.
// It returns a pointer to the newly created controller with an empty clients map.
func NewRealTimeController() *RealTimeController {
	return &RealTimeController{
		clients: make(map[chan string]bool), // Initialize an empty map for client channels.
	}
}

// RegisterRoutes registers the SSE endpoints with the given router group.
// It adds two routes: GET /updates for opening an SSE connection, and
// POST /broadcast to send messages to all connected clients.
func (rtc *RealTimeController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/updates", rtc.HandleSSE)           // Clients use this route to open a persistent connection.
	rg.POST("/broadcast", rtc.BroadcastMessage) // This route is used to send a message to all clients.
}

// HandleSSE handles GET /api/updates.
// It establishes a Server-Sent Events connection and streams messages to the client.
func (rtc *RealTimeController) HandleSSE(c *gin.Context) {
	// Create a new channel to send messages to this client.
	messageChan := make(chan string)

	// Lock the clients map to add the new channel.
	rtc.mu.Lock()
	rtc.clients[messageChan] = true // Add the channel to the set of active clients.
	rtc.mu.Unlock()

	// When the connection is closed, remove the client channel.
	defer func() {
		rtc.mu.Lock()
		delete(rtc.clients, messageChan) // Remove channel from the map.
		rtc.mu.Unlock()
		close(messageChan) // Close the channel to free resources.
	}()

	// Use Gin's streaming functionality to continuously send messages.
	c.Stream(func(w io.Writer) bool {
		// Wait for a message on the client's channel.
		msg, ok := <-messageChan
		if !ok {
			// If the channel is closed, stop streaming.
			return false
		}
		// Send an SSE event with the event type "update" and the message.
		c.SSEvent("update", msg)
		return true // Continue streaming.
	})
}

// BroadcastMessage handles POST /api/broadcast.
// It reads a JSON payload with a "message" field and sends this message
// to all connected SSE clients.
func (rtc *RealTimeController) BroadcastMessage(c *gin.Context) {
	// Define a payload structure to bind the incoming JSON.
	var payload struct {
		Message string `json:"message"` // Message to broadcast.
	}
	// Bind the JSON payload from the request.
	if err := c.ShouldBindJSON(&payload); err != nil {
		// If binding fails, return a 400 Bad Request with the error.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lock the clients map to safely iterate over it.
	rtc.mu.Lock()
	defer rtc.mu.Unlock()

	// Send the message to every connected client channel.
	for ch := range rtc.clients {
		ch <- payload.Message
	}

	// Respond with HTTP 200 and a confirmation message.
	c.JSON(http.StatusOK, gin.H{"status": "broadcasted"})
}
