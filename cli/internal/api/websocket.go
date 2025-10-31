package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type DeploymentUpdate struct {
	DeploymentID string            `json:"deploymentId"`
	ProjectID    string            `json:"projectId"`
	CommitHash   string            `json:"commitHash"`
	Status       string            `json:"status"`
	Namespace    string            `json:"namespace,omitempty"`
	PVC          string            `json:"pvc,omitempty"`
	URL          string            `json:"url,omitempty"`
	Logs         []string          `json:"logs"`
}

type WebSocketClient struct {
	baseURL string
	conn    *websocket.Conn
}

func NewWebSocketClient(baseURL string) *WebSocketClient {
	if baseURL == "" {
		baseURL = DefaultAPIURL
	}
	// Store base URL as-is - we'll convert to WebSocket in Connect()
	return &WebSocketClient{baseURL: baseURL}
}

func (c *WebSocketClient) Connect(deploymentID string) error {
	// Convert HTTP URL to WebSocket URL
	wsBase := c.baseURL
	if strings.HasPrefix(wsBase, "http://") {
		wsBase = "ws://" + wsBase[7:]
	} else if strings.HasPrefix(wsBase, "https://") {
		wsBase = "wss://" + wsBase[8:]
	} else if !strings.HasPrefix(wsBase, "ws://") && !strings.HasPrefix(wsBase, "wss://") {
		// No protocol specified, assume ws://
		wsBase = "ws://" + wsBase
	}
	
	// Build full WebSocket URL
	url := fmt.Sprintf("%s/ws?deploymentId=%s", wsBase, deploymentID)
	
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket at %s: %w", url, err)
	}

	c.conn = conn
	return nil
}

func (c *WebSocketClient) StreamUpdates(callback func(*DeploymentUpdate) error) error {
	if c.conn == nil {
		return fmt.Errorf("not connected - call Connect() first")
	}

	defer c.conn.Close()

	// Set read deadline to detect connection issues
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var update DeploymentUpdate
		err := c.conn.ReadJSON(&update)
		if err != nil {
			// Check if it's a close error
			if closeErr, ok := err.(*websocket.CloseError); ok {
				// Normal closure or going away - server closed connection normally
				if closeErr.Code == websocket.CloseNormalClosure || closeErr.Code == websocket.CloseGoingAway {
					return nil // Return nil to indicate normal completion
				}
				// Unexpected close error
				return fmt.Errorf("WebSocket closed with code %d: %w", closeErr.Code, err)
			}
			// Handle read deadline exceeded (timeout)
			if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() {
				return fmt.Errorf("WebSocket read timeout - connection may be stale")
			}
			// Other errors (network errors, JSON decode errors, etc.)
			return fmt.Errorf("failed to read WebSocket message: %w", err)
		}

		// Call callback with update
		if err := callback(&update); err != nil {
			return err
		}

		// Exit if deployment is complete or failed
		if update.Status == "complete" || update.Status == "failed" {
			return nil
		}

		// Reset read deadline
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	}
}

func (c *WebSocketClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

