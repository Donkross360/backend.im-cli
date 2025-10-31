package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const DefaultAPIURL = "http://localhost:8080"

type Client struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

func NewClient() *Client {
	apiURL := os.Getenv("BACKEND_IM_API_URL")
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}
	
	// Debug: log the API URL being used (can remove later)
	// fmt.Fprintf(os.Stderr, "DEBUG: Using API URL: %s\n", apiURL)

	return &Client{
		baseURL: apiURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) SetAuthToken(token string) {
	c.authToken = token
}

func (c *Client) GenerateCode(prompt string) (map[string]string, error) {
	reqBody := map[string]string{
		"prompt": prompt,
	}

	var response struct {
		Files map[string]string `json:"files"`
	}

	err := c.post("/api/generate", reqBody, &response)
	if err != nil {
		return nil, err
	}

	return response.Files, nil
}

func (c *Client) Deploy(files map[string]string, projectID string) (*DeployResponse, error) {
	reqBody := map[string]interface{}{
		"files":     files,
		"projectId": projectID,
	}

	var response DeployResponse
	err := c.post("/api/deploy", reqBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetStatus(deploymentID string) (*StatusResponse, error) {
	var response StatusResponse
	err := c.get(fmt.Sprintf("/api/status/%s", deploymentID), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) VerifyAuth() (*AuthVerifyResponse, error) {
	var response AuthVerifyResponse
	err := c.get("/api/auth/verify", &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) CommitChanges(files map[string]string, projectID, message string) (*CommitResponse, error) {
	reqBody := map[string]interface{}{
		"files":     files,
		"projectId": projectID,
		"message":   message,
	}

	var response CommitResponse
	err := c.post("/api/commit", reqBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) post(path string, body interface{}, response interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) get(path string, response interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Response types
type DeployResponse struct {
	DeploymentID string `json:"deploymentId"`
	ProjectID    string `json:"projectId"`
	CommitHash   string `json:"commitHash"`
	Status       string `json:"status"`
	WebSocketURL string `json:"websocketUrl,omitempty"`
}

type StatusResponse struct {
	ID         string   `json:"id"`
	ProjectID  string   `json:"projectId"`
	CommitHash string   `json:"commitHash"`
	Status     string   `json:"status"`
	URL        string   `json:"url"`
	Logs       []string `json:"logs"`
}

type AuthVerifyResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

type CommitResponse struct {
	CommitHash string `json:"commitHash"`
	ProjectID  string `json:"projectId"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

