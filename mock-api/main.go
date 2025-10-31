package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	http.HandleFunc("/api/generate", mockGenerate)
	http.HandleFunc("/api/deploy", mockDeploy)
	http.HandleFunc("/api/commit", mockCommit)
	http.HandleFunc("/api/auth/callback", mockAuthCallback)
	http.HandleFunc("/api/auth/verify", mockVerifyAuth)
	http.HandleFunc("/api/status/", mockStatus)
	http.HandleFunc("/ws", mockWebSocket)

	fmt.Println("🚀 Mock Backend.im API running on :8080")
	fmt.Println("📡 WebSocket endpoint: ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// POST /api/generate - Returns mock FastAPI code
func mockGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Simulate API delay (Backend.im calling OpenAI/Claude)
	time.Sleep(2 * time.Second)

	// Return mock FastAPI code (simulating what Backend.im would return)
	response := map[string]interface{}{
		"files": map[string]string{
			"main.py": fmt.Sprintf(`from fastapi import FastAPI

app = FastAPI()

@app.get("/")
def read_root():
    return {"message": "Generated from: %s"}

@app.get("/health")
def health():
    return {"status": "healthy"}`, req.Prompt),
			"models.py": `from sqlalchemy import Column, Integer, String
from database import Base

class User(Base):
    __tablename__ = "users"
    id = Column(Integer, primary_key=True)
    name = Column(String(50))`,
			"requirements.txt": "fastapi==0.104.1\nuvicorn==0.24.0\nsqlalchemy==2.0.23",
			"schema.sql": "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(50));",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /api/deploy - Returns deployment ID, project ID, and commit hash
func mockDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Files     map[string]string `json:"files"`
		ProjectID string           `json:"projectId"` // Unique project ID (includes user ID)
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	deploymentID := uuid.New().String()
	commitHash := generateCommitHash(req.Files)

	// Project ID should be provided by client (created when user creates project)
	projectID := req.ProjectID
	if projectID == "" {
		projectID = "proj-" + uuid.New().String()[:8] // Mock fallback
	}

	// Track deployment start time for status polling
	deploymentStartTimes[deploymentID] = time.Now()

	response := map[string]interface{}{
		"deploymentId": deploymentID,
		"projectId":   projectID,   // Used with commit hash for namespace: {projectId}-{commitHash}
		"commitHash":   commitHash,  // Combined with project ID for unique namespace/PVC
		"status":       "queued",
		"websocketUrl": fmt.Sprintf("ws://localhost:8080/ws?deploymentId=%s", deploymentID),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /api/commit - Commits local changes to Backend.im/Gitea
func mockCommit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Files     map[string]string `json:"files"`
		ProjectID string            `json:"projectId"`
		Message   string            `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	commitHash := generateCommitHash(req.Files)

	response := map[string]interface{}{
		"commitHash": commitHash,
		"projectId":  req.ProjectID,
		"status":     "committed",
		"message":    req.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Simple in-memory deployment tracking
var deploymentStatuses = make(map[string]map[string]interface{})
var deploymentStartTimes = make(map[string]time.Time)

// GET /api/status/{deploymentId} - Returns current deployment status
func mockStatus(w http.ResponseWriter, r *http.Request) {
	deploymentID := r.URL.Path[len("/api/status/"):]
	if deploymentID == "" {
		http.Error(w, "Deployment ID required", http.StatusBadRequest)
		return
	}

	// Check if we have tracked status for this deployment
	if status, exists := deploymentStatuses[deploymentID]; exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
		return
	}

	// Check if deployment just started (for polling scenarios)
	startTime, exists := deploymentStartTimes[deploymentID]
	if !exists {
		// New deployment - mark start time
		deploymentStartTimes[deploymentID] = time.Now()
		startTime = time.Now()
	}

	elapsed := time.Since(startTime)
	commitHash := "a1b2c3d4e5f6" // Mock commit hash

	// Simulate deployment stages based on elapsed time
	var status string
	var logs []string
	var url string

	if elapsed < 3*time.Second {
		status = "queued"
		logs = []string{"📦 Deployment queued..."}
	} else if elapsed < 6*time.Second {
		status = "committing"
		logs = []string{"📦 Committing files to repository..."}
	} else if elapsed < 9*time.Second {
		status = "creating_namespace"
		logs = []string{"🏗️ Creating Kubernetes namespace..."}
	} else if elapsed < 12*time.Second {
		status = "creating_pvc"
		logs = []string{"💾 Creating Persistent Volume Claim..."}
	} else if elapsed < 15*time.Second {
		status = "building"
		logs = []string{"🔨 Building container image..."}
	} else if elapsed < 18*time.Second {
		status = "deploying"
		logs = []string{"🚀 Deploying to Kubernetes cluster..."}
	} else {
		// Complete after ~18 seconds
		status = "complete"
		url = fmt.Sprintf("https://%s.backend.im", deploymentID[:12])
		logs = []string{"✅ Deployment complete!"}
		
		// Store completed status
		deploymentStatuses[deploymentID] = map[string]interface{}{
			"id":         deploymentID,
			"projectId":  "user123-myproject",
			"commitHash": commitHash,
			"status":     status,
			"url":        url,
			"logs":       logs,
		}
	}

	response := map[string]interface{}{
		"id":         deploymentID,
		"projectId":  "user123-myproject",
		"commitHash": commitHash,
		"status":     status,
		"logs":       logs,
	}

	if url != "" {
		response["url"] = url
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// WebSocket /ws?deploymentId={id} - Streams deployment updates
// Orchestrator needs project ID + commit hash to create unique namespace/PVC
func mockWebSocket(w http.ResponseWriter, r *http.Request) {
	deploymentID := r.URL.Query().Get("deploymentId")
	if deploymentID == "" {
		http.Error(w, "deploymentId required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// In real implementation, these would come from deployment record
	projectID := "user123-myproject" // Unique project ID (includes user ID)
	commitHash := "a1b2c3d4e5f6"     // Commit hash for this deployment

	// Stream deployment stages - project ID + commit hash required for orchestrator
	stages := []map[string]interface{}{
		{
			"deploymentId": deploymentID,
			"projectId":   projectID,   // Required: identifies the project
			"commitHash":   commitHash,  // Required: orchestrator uses this to pull from Gitea
			"status":       "committing",
			"logs":         []string{"📦 Committing files to repository..."},
		},
		{
			"deploymentId": deploymentID,
			"projectId":   projectID,   // Orchestrator combines projectId + commitHash for namespace
			"commitHash":   commitHash,  // Namespace format: {projectId}-{commitHash}
			"status":       "creating_namespace",
			"logs":         []string{"🏗️ Creating Kubernetes namespace..."},
		},
		{
			"deploymentId": deploymentID,
			"projectId":   projectID,   // Same project ID, different commit = update deployment
			"commitHash":   commitHash,  // PVC format: {projectId}-{commitHash}
			"status":       "creating_pvc",
			"logs":         []string{"💾 Creating Persistent Volume Claim..."},
		},
		{
			"deploymentId": deploymentID,
			"projectId":   projectID,
			"commitHash":   commitHash,
			"status":       "deploying",
			"logs":         []string{"🚀 Deploying application..."},
		},
		{
			"deploymentId": deploymentID,
			"projectId":   projectID,
			"commitHash":   commitHash,
			"status":       "complete",
			"url":          fmt.Sprintf("https://%s.backend.im", deploymentID[:12]), // Mock URL - for testing only
			"logs":         []string{"✅ Deployment complete!"},
		},
	}

	for _, stage := range stages {
		if err := conn.WriteJSON(stage); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
		time.Sleep(2 * time.Second)
	}
}

func generateCommitHash(files map[string]string) string {
	hasher := sha256.New()
	for filename, content := range files {
		hasher.Write([]byte(filename + content))
	}
	return hex.EncodeToString(hasher.Sum(nil))[:12]
}

func mockAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Return mock JWT token
	response := map[string]interface{}{
		"access_token": "mock_token_" + uuid.New().String(),
		"token_type":   "Bearer",
		"expires_in":   3600,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func mockVerifyAuth(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}

	// Simple token verification (in real app, would validate JWT)
	response := map[string]interface{}{
		"valid":  true,
		"userId": "user123",
		"email":  "user@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

