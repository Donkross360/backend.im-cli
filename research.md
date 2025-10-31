# 🚀 Backend.im CLI Integration Research Report
## DevOps Research Task – Infrastructure Setup & CLI Flow for Backend.im

**Research Objective**: Design a CLI-based deployment workflow that enables developers to deploy backend code directly to Backend.im via Claude Code CLI and other AI tools, using mostly open-source tools with minimal configuration.

---

## 📋 Executive Summary

This research proposes a CLI-based deployment solution that integrates with the existing Backend.im platform. The solution leverages a Go-based CLI tool that acts as a thin wrapper around Backend.im's existing APIs, providing developers with a seamless CLI-first deployment experience while maintaining the security and reliability of the existing platform.

**Key Benefits:**
- ✅ **Zero Infrastructure Changes**: Leverages existing Backend.im platform
- ✅ **Developer-Focused**: CLI-first approach for backend developers
- ✅ **Cost-Effective**: Minimal additional infrastructure required
- ✅ **Secure**: Reuses existing authentication and security systems
- ✅ **Scalable**: Go-based CLI handles enterprise-scale usage

---

## 🎯 Problem Analysis

### Current Backend.im Workflow
The existing Backend.im platform provides a web-based interface where:
1. Users enter prompts in a small input box
2. Backend sends prompts to OpenAI/Claude for code generation
3. Generated code is displayed in Monaco Editor (web-based code editor) for editing
4. Users can modify code before saving to Gitea
5. Deployment orchestrator creates K8s namespace and PVC
6. Application is deployed and URL is returned

### Developer Pain Points
- **Context Switching**: Developers must switch from CLI to web browser
- **Workflow Disruption**: Breaks the natural CLI-based development flow
- **Limited Automation**: No way to script or automate deployments
- **Review Process**: Difficult to review large amounts of generated code in web UI

### Solution Requirements
- **CLI-First Experience**: Developers work primarily from command line
- **Code Review**: Ability to review and edit generated code before deployment
- **One-Command Deployment**: Simple, automated deployment process
- **Local Editing**: Ability to edit generated code locally with any editor
- **Minimal Infrastructure**: No changes to existing Backend.im platform

---

## 🏗️ Proposed Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    Developer Environment                        │
├─────────────────────────────────────────────────────────────────┤
│  Go CLI Tool  →  Local Files (for editing)  →  Optional Editor │
│  (Backend.im handles Git commits - same repos as UI)           │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Backend.im Platform (Existing)             │
├─────────────────────────────────────────────────────────────────┤
│  API Gateway  →  Auth Service  →  Git Server (Gitea)         │
│       │              │                    │                    │
│       ▼              ▼                    ▼                    │
│  WebSocket     Redis Queue    Deployment Orchestrator          │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                 Kubernetes Infrastructure (Existing)           │
├─────────────────────────────────────────────────────────────────┤
│  Namespace Management  →  PVC Provisioning  →  App Deployment  │
│         │                       │                    │         │
│         ▼                       ▼                    ▼         │
│  Service Mesh        Persistent Storage    Load Balancer       │
└─────────────────────────────────────────────────────────────────┘
```

### Core Components

#### 1. **Go CLI Tool (NEW)**
- **Purpose**: Primary interface for developers
- **Technology**: Go with Cobra CLI framework
- **Features**: Code generation via Backend.im API, file download, deployment management

#### 2. **Local File Management (NEW)**
- **Purpose**: Download local copy of code for editing (Backend.im is source of truth)
- **Technology**: File system operations
- **Features**: File download, optional editor integration (VS Code, vim, etc.)
- **Note**: No local Git management - Backend.im handles all Git commits automatically

#### 3. **Backend.im Platform (EXISTING)**
- **Purpose**: Core deployment infrastructure
- **Components**: API Gateway, Auth Service, Gitea, Kubernetes
- **Status**: No changes required

---

## 🛠️ Technology Stack & Tool Selection

### CLI Development Stack

| Component | Tool | Reasoning |
|-----------|------|-----------|
| **CLI Framework** | Go + Cobra | Enterprise-grade, zero dependencies, fast execution |
| **API Client** | Go net/http | Built-in HTTP client, no external dependencies |
| **Configuration** | Viper | Go-native configuration management |
| **Authentication** | OAuth 2.0 | Reuse existing Backend.im auth system |
| **File Operations** | Go os/io | Built-in file system operations |

### Editor Integration Approach

| Approach | Pros | Cons | Recommendation |
|----------|------|------|----------------|
| **No Editor Integration** | Simple, no dependencies | User must manually open files | ✅ **Recommended for v1** |
| **Optional Editor Flag** | Flexible, user choice | Requires editor detection | ✅ **Consider for v1** |
| **VS Code Extension** | Rich integration | Complex development | Consider for v2 |

### Existing Infrastructure (No Changes Required)

| Component | Status | Purpose |
|-----------|--------|---------|
| **Backend.im API** | Existing | Handles deployment requests |
| **Gitea Server** | Existing | Git repository hosting |
| **Kubernetes Cluster** | Existing | Container orchestration |
| **WebSocket Service** | Existing | Real-time deployment updates |
| **User Authentication** | Existing | OAuth and user management |
| **Deployment Orchestrator** | Existing | Handles K8s namespace and PVC creation |

---

## 🔄 Deployment Sequence Flow

### Phase 1: Code Generation (NEW - CLI Tool)

```
Developer → Go CLI → Backend.im API → Generated Code (FastAPI + DB Schemas) → Download to Local Files
```

**Steps:**
1. Developer runs: `backend-im generate "Create a REST API for user management"`
2. CLI sends prompt to Backend.im API (`POST /api/generate`)
3. Backend.im generates FastAPI code + DB schemas (as it does in web UI)
4. **Backend.im auto-commits generated code to Gitea** (same as UI behavior)
5. CLI receives generated code files (main.py, models.py, requirements.txt, etc.)
6. CLI downloads files to local project directory (for editing)
7. Optional: CLI can open editor (`--editor vscode` or `--editor code`) if specified

**Note**: Generated code is already committed to Gitea by Backend.im - CLI just downloads a local copy for editing

### Phase 2: Code Review (NEW - Local Editing)

```
Local Editor → Code Editing → User Review → Save/Deploy Decision
```

**Steps:**
1. Generated code files are saved locally in project directory (local copy for editing)
2. User edits code using their preferred editor (VS Code, vim, nano, etc.)
3. User can use extensions, debugging, and all local dev tools
4. User saves changes to local files
5. When ready, user runs `backend-im deploy` to upload changes
6. **Backend.im commits changes to Gitea** (same as UI save button)
7. **Same project visible in UI** - user can switch between CLI and web UI seamlessly

### Phase 3: Deployment (EXISTING - Backend.im Platform)

```
CLI → Backend.im API → Gitea → Kubernetes → WebSocket Updates
```

**Steps:**
1. Developer runs: `backend-im deploy` (or `backend-im deploy --watch` for live updates)
2. CLI reads local project files (FastAPI code, DB schemas, requirements.txt)
3. CLI uploads code to Backend.im API (`POST /api/deploy`)
4. **Backend.im API commits code to Gitea** (same as UI save button)
5. **Same Git repo** - projects visible in both CLI and UI
6. Existing deployment orchestrator watches Gitea, creates K8s namespace and PVC
7. Application is deployed and URL is returned
8. CLI receives real-time updates via WebSocket (if `--watch` flag used)

**Important**: Backend.im already handles Git commits automatically:
- Code generation → Auto-commits to Gitea
- UI save button → Commits changes to Gitea
- CLI upload → Should commit to same Gitea repo (consistent with UI)
- **No local Git management needed** - Backend.im is the source of truth

---

## 💻 Local Setup Flow

### Prerequisites Installation

```bash
# 1. Install Go (required for CLI tool)
curl -fsSL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz | sudo tar -xzC /usr/local

# 2. Git not required - Backend.im handles all Git operations automatically
```

### Backend.im CLI Setup

```bash
# 1. Install Backend.im CLI tool
go install github.com/backend-im/cli@latest

# 2. Authenticate with existing Backend.im platform
backend-im auth

# 3. Generate code (Backend.im commits to Gitea automatically, CLI downloads local copy)
backend-im generate "Create a REST API for user management" --project my-api

# 4. Edit code locally (optional - can also edit in UI)
code .  # or vim, nano, etc.

# 5. Deploy changes (Backend.im commits to same Gitea repo as UI)
backend-im deploy

# 6. View project in UI - same project visible in both CLI and web interface
# User can switch between CLI and UI seamlessly
```

### No Additional Infrastructure Required

**Important**: No additional infrastructure setup is required because:
- Backend.im platform already exists and handles all deployment infrastructure
- **Backend.im already handles Git commits** - CLI doesn't need local Git management
- **Same Gitea repos** - projects created via CLI are visible in UI and vice versa
- CLI tool only needs to communicate with existing API endpoints
- Code generation happens on Backend.im platform (no local AI services needed)
- No new servers, databases, or Kubernetes clusters needed

---

## 🔧 Minimal Custom Code Requirements

### 1. Go CLI Tool (NEW - Core Implementation)

**File: `cmd/backend-im/main.go`**
```go
package main

import (
    "github.com/spf13/cobra"
    "github.com/backend-im/cli/internal/commands"
)

func main() {
    var rootCmd = &cobra.Command{
        Use:   "backend-im",
        Short: "Backend.im CLI for seamless deployment",
        Long:  "A CLI tool for deploying backend code to Backend.im platform",
    }
    
    rootCmd.AddCommand(commands.NewGenerateCommand())
    rootCmd.AddCommand(commands.NewDeployCommand())
    rootCmd.AddCommand(commands.NewAuthCommand())
    rootCmd.Execute()
}
```

**File: `internal/commands/generate.go`**
```go
func NewGenerateCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "generate [prompt]",
        Short: "Generate backend code from prompt via Backend.im API",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            prompt := args[0]
            editor, _ := cmd.Flags().GetString("editor")
            outputDir, _ := cmd.Flags().GetString("output")
            
            // Call Backend.im API to generate code
            apiClient := api.NewClient()
            files, err := apiClient.GenerateCode(prompt)
            if err != nil {
                return fmt.Errorf("failed to generate code: %w", err)
            }
            
            // Download files to local directory
            projectDir := outputDir
            if projectDir == "" {
                projectDir = fmt.Sprintf("backend-%d", time.Now().Unix())
            }
            
            err = downloadFiles(files, projectDir)
            if err != nil {
                return fmt.Errorf("failed to download files: %w", err)
            }
            
            fmt.Printf("✅ Code generated successfully in ./%s\n", projectDir)
            
            // Optionally open editor
            if editor != "" {
                return openEditor(editor, projectDir)
            }
            
            return nil
        },
    }
    
    cmd.Flags().StringP("editor", "e", "", "Open generated code in editor (vscode, code, vim, etc.)")
    cmd.Flags().StringP("output", "o", "", "Output directory (default: auto-generated name)")
    return cmd
}
```

**File: `internal/commands/deploy.go`**
```go
func NewDeployCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "deploy",
        Short: "Deploy local code to Backend.im",
        RunE: func(cmd *cobra.Command, args []string) error {
            watch, _ := cmd.Flags().GetBool("watch")
            projectDir, _ := cmd.Flags().GetString("dir")
            
            if projectDir == "" {
                projectDir = "."
            }
            
            // Read local project files
            files, err := readProjectFiles(projectDir)
            if err != nil {
                return fmt.Errorf("failed to read project files: %w", err)
            }
            
            // Upload and deploy via Backend.im API
            apiClient := api.NewClient()
            result, err := apiClient.Deploy(files)
            if err != nil {
                return fmt.Errorf("deployment failed: %w", err)
            }
            
            fmt.Printf("🚀 Deployment started: %s\n", result.URL)
            fmt.Printf("📋 Deployment ID: %s\n", result.DeploymentID)
            
            // Watch deployment status if requested
            if watch {
                return watchDeployment(apiClient, result.DeploymentID)
            }
            
            return nil
        },
    }
    
    cmd.Flags().BoolP("watch", "w", false, "Watch deployment progress in real-time")
    cmd.Flags().StringP("dir", "d", "", "Project directory (default: current directory)")
    return cmd
}
```

### 2. File Download & Project Management (NEW - Minimal Implementation)

**File: `internal/files/project.go`**
```go
package files

import (
    "archive/tar"
    "compress/gzip"
    "io"
    "os"
    "path/filepath"
)

// DownloadFiles extracts generated code files from API response to local directory
func DownloadFiles(files map[string][]byte, projectDir string) error {
    err := os.MkdirAll(projectDir, 0755)
    if err != nil {
        return err
    }
    
    // Write each file to project directory
    for filename, content := range files {
        filePath := filepath.Join(projectDir, filename)
        
        // Create directory if needed
        dir := filepath.Dir(filePath)
        os.MkdirAll(dir, 0755)
        
        err := os.WriteFile(filePath, content, 0644)
        if err != nil {
            return err
        }
    }
    
    return nil
}

// ReadProjectFiles reads all project files from local directory
func ReadProjectFiles(projectDir string) (map[string][]byte, error) {
    files := make(map[string][]byte)
    
    err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return nil
        }
        
        // Skip hidden files and common ignore patterns
        if filepath.Base(path)[0] == '.' {
            return nil
        }
        
        content, err := os.ReadFile(path)
        if err != nil {
            return err
        }
        
        relPath, _ := filepath.Rel(projectDir, path)
        files[relPath] = content
        
        return nil
    })
    
    return files, err
}
```

**File: `internal/editor/integration.go`**
```go
package editor

import (
    "os/exec"
    "os"
)

// OpenEditor opens the specified editor in the given directory
func OpenEditor(editorName, projectDir string) error {
    var cmd *exec.Cmd
    
    switch editorName {
    case "vscode", "code":
        cmd = exec.Command("code", projectDir)
    case "vim":
        cmd = exec.Command("vim", projectDir)
    case "nano":
        cmd = exec.Command("nano", projectDir)
    default:
        // Try to run as-is (user might specify full path)
        cmd = exec.Command(editorName, projectDir)
    }
    
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Start() // Don't wait, let editor run independently
}
```

### 3. API Client (NEW - Backend.im Integration)

**File: `internal/api/client.go`**
```go
package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type Client struct {
    baseURL    string
    authToken  string
    httpClient *http.Client
}

// GenerateCode calls Backend.im API to generate FastAPI code + DB schemas
func (c *Client) GenerateCode(prompt string) (map[string][]byte, error) {
    req := &GenerateRequest{
        Prompt: prompt,
    }
    
    resp, err := c.post("/api/generate", req)
    if err != nil {
        return nil, err
    }
    
    // Parse response - Backend.im returns files as map or archive
    var result GenerateResponse
    json.Unmarshal(resp, &result)
    
    // Download files if URLs provided, or extract from response
    files := make(map[string][]byte)
    for filename, content := range result.Files {
        files[filename] = []byte(content)
    }
    
    return files, nil
}

// Deploy uploads local files and triggers deployment
func (c *Client) Deploy(files map[string][]byte) (*DeployResponse, error) {
    req := &DeployRequest{
        Files: files,
    }
    
    resp, err := c.post("/api/deploy", req)
    if err != nil {
        return nil, err
    }
    
    var result DeployResponse
    json.Unmarshal(resp, &result)
    return &result, nil
}

// WatchDeployment streams deployment status via WebSocket
func (c *Client) WatchDeployment(deploymentID string) error {
    // Connect to WebSocket endpoint
    // Stream status updates to terminal
    return nil
}
```

### 4. API Contract Discovery & Mocking Strategy

Since the Backend.im web UI exists, we can reverse-engineer the API contracts and create a mock API for parallel development:

**Step 1: Reverse-Engineer Existing Web UI**
```bash
# Enable browser DevTools and network monitoring
# 1. Open Backend.im web UI
# 2. Open Chrome DevTools → Network tab
# 3. Trigger deployment flow
# 4. Capture all API requests/responses
```

**Key API Contracts to Discover:**
- `/api/generate` - Code generation endpoint (prompt → FastAPI + DB schemas)
- `/api/auth/login` - OAuth callback handling
- `/api/deploy` - Deployment submission (files → deployment)
- `/api/status/{deploymentId}` - Deployment status
- `/ws` - WebSocket endpoint for real-time updates

**Step 2: Create Mock API Server**

**File: `mock-api/main.go`** (Simple mock for testing)
```go
package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/google/uuid"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
    http.HandleFunc("/api/generate", mockGenerate)
    http.HandleFunc("/api/deploy", mockDeploy)
    http.HandleFunc("/api/auth/callback", mockAuthCallback)
    http.HandleFunc("/api/auth/verify", mockVerifyAuth)
    http.HandleFunc("/ws", mockWebSocket)
    
    fmt.Println("Mock Backend.im API running on :8080")
    http.ListenAndServe(":8080", nil)
}

// POST /api/generate - Returns mock FastAPI code
func mockGenerate(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Prompt string `json:"prompt"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    
    time.Sleep(2 * time.Second) // Simulate API delay
    
    json.NewEncoder(w).Encode(map[string]interface{}{
        "files": map[string]string{
            "main.py": `from fastapi import FastAPI
app = FastAPI()
@app.get("/")
def read_root():
    return {"message": "Hello World"}`,
            "requirements.txt": "fastapi==0.104.1\nuvicorn==0.24.0",
        },
    })
}

// POST /api/deploy - Returns deployment ID, project ID, and commit hash
func mockDeploy(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Files     map[string]string `json:"files"`
        ProjectID string           `json:"projectId"` // Unique project ID (includes user ID)
    }
    json.NewDecoder(r.Body).Decode(&req)
    
    deploymentID := uuid.New().String()
    commitHash := generateCommitHash(req.Files)
    
    // Project ID should be provided by client (created when user creates project)
    // Format: userID-projectName or similar unique identifier
    projectID := req.ProjectID
    if projectID == "" {
        projectID = "proj-" + uuid.New().String()[:8] // Mock fallback
    }
    
    json.NewEncoder(w).Encode(map[string]interface{}{
        "deploymentId": deploymentID,
        "projectId":   projectID,   // Used with commit hash for namespace: {projectId}-{commitHash}
        "commitHash":   commitHash,  // Combined with project ID for unique namespace/PVC
        "status":       "queued",
    })
}

// WebSocket /ws?deploymentId={id} - Streams deployment updates
// Orchestrator needs project ID + commit hash to create unique namespace/PVC
func mockWebSocket(w http.ResponseWriter, r *http.Request) {
    deploymentID := r.URL.Query().Get("deploymentId")
    
    // In real implementation, these would come from deployment record
    projectID := "user123-myproject" // Unique project ID (includes user ID)
    commitHash := "a1b2c3d4e5f6"     // Commit hash for this deployment
    
    conn, _ := upgrader.Upgrade(w, r, nil)
    defer conn.Close()
    
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
            "url":          fmt.Sprintf("https://%s.backend.im", deploymentID[:12]),
            "logs":         []string{"✅ Deployment complete!"},
        },
    }
    
    for _, stage := range stages {
        conn.WriteJSON(stage)
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
    json.NewEncoder(w).Encode(map[string]interface{}{
        "access_token": "mock_token_" + uuid.New().String(),
        "token_type":   "Bearer",
    })
}

func mockVerifyAuth(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]bool{"valid": true})
}
```

**Key Features of Mock API:**

1. **Simple Simulation**:
   - Generates commit hash from file contents (deterministic SHA256)
   - Returns mock FastAPI code (simulates OpenAI/Claude)
   - Simple WebSocket streaming with commit hash
   - Basic OAuth mock (no real Google calls)

2. **WebSocket Streaming**:
   - **Project ID + Commit hash are required** - orchestrator uses them to:
     - Pull the correct commit from Gitea (using commit hash)
     - Create unique namespace names: `{projectId}-{commitHash}`
     - Create unique PVC names: `{projectId}-{commitHash}`
     - If namespace exists for same projectId, orchestrator deletes old one first
   - Streams status updates and logs
   - Streams final URL when deployment completes
   - Note: Project ID is created when user creates project (includes user ID for uniqueness)
   - Note: Same project, different commit = same projectId, different commitHash = namespace gets recreated

3. **No Complex State**:
   - Simple in-memory responses
   - No database or complex state management
   - Just enough to test CLI integration

**Benefits of Mocking Approach:**
- ✅ Parallel development (CLI + Backend.im platform)
- ✅ Contract-driven development
- ✅ Integration testing without production risk
- ✅ Documentation of expected API behavior
- ✅ Can be used for demos while platform is being built
- ✅ Tests complete deployment flow including orchestrator metadata
- ✅ Validates WebSocket streaming with commit hash, namespace, PVC, and URL

**Step 3: API Contract Documentation**

After reverse-engineering, document the contracts:

**File: `API_CONTRACTS.md`**
```markdown
## Backend.im API Contracts

### POST /api/generate
**Purpose**: Generate FastAPI code + DB schemas from prompt (simulates Backend.im calling OpenAI/Claude)

Request:
```json
{
  "prompt": "Create a REST API for user management"
}
```

Response:
```json
{
  "files": {
    "main.py": "# FastAPI code...",
    "models.py": "# SQLAlchemy models...",
    "requirements.txt": "fastapi==0.104.1\n...",
    "schema.sql": "CREATE TABLE users..."
  }
}
```

### POST /api/deploy
**Purpose**: Submit code files for deployment, returns deployment ID, project ID, and commit hash

Request:
```json
{
  "files": {
    "main.py": "# FastAPI code...",
    "models.py": "# SQLAlchemy models...",
    "requirements.txt": "fastapi==0.104.1\n...",
    "schema.sql": "CREATE TABLE users..."
  },
  "projectId": "user123-myproject"  // Unique project ID (created when user creates project, includes user ID)
}
```

Response:
```json
{
  "deploymentId": "550e8400-e29b-41d4-a716-446655440000",
  "projectId": "user123-myproject",  // Used with commit hash for namespace: {projectId}-{commitHash}
  "commitHash": "a1b2c3d4e5f6",       // Combined with project ID for unique namespace/PVC
  "status": "queued",
  "websocketUrl": "ws://backend.im/ws?deploymentId=550e8400-e29b-41d4-a716-446655440000"
}
```

**Project ID Format**:
- Created when user first creates a project (after login)
- Includes user ID to ensure uniqueness: `{userId}-{projectName}` or similar
- Remains constant for that project across all deployments
- Example: `user123-myapi`, `user456-blog-backend`

### GET /api/status/{deploymentId}
**Purpose**: Get current deployment status (optional polling fallback)

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "projectId": "user123-myproject",
  "commitHash": "a1b2c3d4e5f6",
  "status": "deploying",
  "url": "https://a1b2c3d4e5f6.backend.im",
  "logs": ["📦 Committing files...", "🏗️ Creating namespace..."]
}
```

**Note**: Namespace and PVC are created by the orchestrator using `{projectId}-{commitHash}` format, not returned in API responses.

### WebSocket /ws?deploymentId={id}
**Purpose**: Stream real-time deployment progress with project ID + commit hash (required for orchestrator)

**Connection**: `ws://backend.im/ws?deploymentId=550e8400-e29b-41d4-a716-446655440000`

**Message Format** (sent repeatedly as deployment progresses):
```json
{
  "deploymentId": "550e8400-e29b-41d4-a716-446655440000",
  "projectId": "user123-myproject",  // REQUIRED: Identifies the project (includes user ID)
  "commitHash": "a1b2c3d4e5f6",      // REQUIRED: Orchestrator uses this to pull from Gitea
  "status": "creating_namespace",
  "logs": [
    "📦 Committing files to repository (commit: a1b2c3d4e5f6)",
    "🏗️ Creating Kubernetes namespace: user123-myproject-a1b2c3d4e5f6"
  ]
}
```

**Complete Example** (when deployment finishes):
```json
{
  "deploymentId": "550e8400-e29b-41d4-a716-446655440000",
  "projectId": "user123-myproject",
  "commitHash": "a1b2c3d4e5f6",
  "status": "complete",
  "url": "https://a1b2c3d4e5f6.backend.im",
  "logs": ["✅ Deployment complete!"]
}
```

**Status Values** (streamed in sequence):
- `queued` - Deployment queued, waiting to start
- `committing` - Committing files to Gitea repository (commit hash generated here)
- `creating_namespace` - Orchestrator creating Kubernetes namespace using `{projectId}-{commitHash}`
- `creating_pvc` - Orchestrator creating PVC using `{projectId}-{commitHash}` for uniqueness
- `building` - Building container image
- `deploying` - Deploying to Kubernetes cluster
- `health_check` - Running health checks
- `complete` - Deployment successful (URL sent here)
- `failed` - Deployment failed (error logs sent)

**Orchestrator Flow** (using project ID + commit hash):
1. **Receive project ID + commit hash** via WebSocket
2. **Check for existing namespace** - if namespace exists for same `projectId`, delete it first
3. **Pull from Gitea** using commit hash: `git checkout {commitHash}`
4. **Create unique namespace** using: `{projectId}-{commitHash}` (e.g., `user123-myproject-a1b2c3d4e5f6`)
5. **Create unique PVC** using: `{projectId}-{commitHash}` (e.g., `user123-myproject-a1b2c3d4e5f6`)
6. **Build & Deploy** container image in namespace with PVC mounted
7. **Health Check** application health
8. **Complete** - Deployment URL generated and sent via WebSocket

**Why project ID + commit hash:**
- **Project ID** (includes user ID) ensures uniqueness across users: `user123-myproject` vs `user456-myproject`
- **Commit hash** enables updates: same project, new commit = `user123-myproject-a1b2c3` → `user123-myproject-f6e5d4`
- **Namespace format**: `{projectId}-{commitHash}` ensures:
  - Different users don't conflict (different project IDs)
  - Same project updates delete old namespace and create new one (orchestrator handles this)
  - No "backend" prefix needed - project ID already includes user identification

**Example Scenarios**:
- User 123, project "myapi", commit `a1b2c3` → namespace: `user123-myapi-a1b2c3`
- User 123, project "myapi", commit `f6e5d4` → namespace: `user123-myapi-f6e5d4` (old one deleted first)
- User 456, project "myapi", commit `a1b2c3` → namespace: `user456-myapi-a1b2c3` (different user, no conflict)

### POST /api/auth/callback
**Purpose**: OAuth callback endpoint (Google OAuth)

Request:
```
GET /api/auth/callback?code=AUTHORIZATION_CODE&state=STATE
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

### GET /api/auth/verify
**Purpose**: Verify authentication token

Request Headers:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Response:
```json
{
  "valid": true,
  "userId": "user-123",
  "email": "user@example.com"
}
```
```

### 5. Mock API Dependencies & Setup

**File: `mock-api/go.mod`**
```go
module mock-api

go 1.21

require (
    github.com/google/uuid v1.5.0
    github.com/gorilla/websocket v1.5.1
)
```

**Running the Mock API:**
```bash
# Start mock API server
cd mock-api
go mod tidy
go run main.go

# Mock API will be available at:
# - HTTP: http://localhost:8080
# - WebSocket: ws://localhost:8080/ws
```

**Note**: This is a simple mock - just enough to test CLI integration. No complex configuration needed.

**Testing the Mock API:**
```bash
# Test code generation
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Create a REST API"}'

# Test deployment
curl -X POST http://localhost:8080/api/deploy \
  -H "Content-Type: application/json" \
  -d '{
    "files": {"main.py": "print(\"hello\")"},
    "projectName": "test",
    "userId": "user-123"
  }'

# Test WebSocket (use wscat or similar tool)
wscat -c "ws://localhost:8080/ws?deploymentId=YOUR_DEPLOYMENT_ID"
```

### 6. No Additional Backend Code Required (Once API Contracts Known)

**Important**: Once API contracts are reverse-engineered and documented:
- CLI can be built against mock API
- Backend.im team can reference contracts for implementation
- Integration becomes straightforward contract mapping
- Mock can be replaced with real API when ready
- Mock API tests complete orchestrator flow including commit hash, namespace, PVC, and URL streaming

---

## ⚠️ Challenges & Mitigations

### Challenge 1: Incomplete API Discovery

**Problem**: Browser DevTools may not capture all API interactions:
- Some requests might be obfuscated or minified
- WebSocket messages might be encoded
- Error responses might not be visible in normal flow
- Authentication tokens might have complex refresh logic

**Mitigation**:
- Use multiple tools: Chrome DevTools, Burp Suite, or mitmproxy
- Test both success and failure scenarios
- Inspect WebSocket frames directly
- Review source code if available (frontend JS bundles)
- Document assumptions and create extensible mock API

**Risk Level**: Medium - Can add 2-3 days to discovery phase

---

### Challenge 2: API Contract Drift

**Problem**: The Backend.im API might change while CLI is being developed:
- New required fields added
- Response format changes
- Authentication flow updates
- Breaking changes in deployment process

**Mitigation**:
- **Version the API contracts** (v1, v2, etc.)
- **Create API contract tests** that validate both mock and real API
- **Establish communication** with Backend.im team for API changes
- **Design CLI for flexibility** - use configuration files for API endpoints
- **Implement feature flags** to handle API version differences

**Risk Level**: Medium - Requires ongoing coordination

---

### Challenge 3: Google OAuth Authentication

**Problem**: CLI needs Google OAuth authentication:
- Headless CLI authentication (no browser)
- Token refresh mechanisms
- Secure token storage
- Google OAuth device flow or authorization code flow

**Mitigation**:
- **Use Google OAuth Device Flow** (OOB flow) - Google supports this natively
- **Alternative: Use installed app flow** - Google provides `gcloud auth login` pattern
- **Use Google's official libraries** (`golang.org/x/oauth2`) - well-tested and documented
- **Cache tokens securely** in local config with encryption
- **Handle token expiration** with automatic refresh

**Risk Level**: Low-Medium - Google OAuth is well-documented for CLI apps

**Recommended Approach**:
```go
import "golang.org/x/oauth2"

// Google OAuth configuration
var googleOAuthConfig = &oauth2.Config{
    ClientID:     clientID,
    ClientSecret: clientSecret,
    Scopes:       []string{"openid", "email", "profile"},
    Endpoint:     google.Endpoint,
}

// For CLI: Use device flow or installed app flow
func (c *Client) Authenticate() error {
    // Option 1: Device Flow (OOB) - user visits URL, enters code
    authURL := googleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
    fmt.Printf("Visit: %s\n", authURL)
    fmt.Print("Enter code: ")
    // ... handle code exchange
    
    // Option 2: Installed app flow (like gcloud)
    // Opens browser automatically, handles callback
    token, err := googleOAuthConfig.Exchange(ctx, code)
    // Store token securely
    return saveToken(token)
}
```

**Google OAuth Resources**:
- Google OAuth 2.0 for installed applications: https://developers.google.com/identity/protocols/oauth2/native-app
- Device flow documentation is available
- Official Go OAuth2 library with Google examples

---

### Challenge 4: WebSocket Real-Time Updates

**Problem**: WebSocket integration for deployment status:
- Connection lifecycle management
- Reconnection logic
- Message parsing and validation
- Handling concurrent deployments

**Mitigation**:
- **Use proven WebSocket library** (gorilla/websocket)
- **Implement exponential backoff** for reconnections
- **Create WebSocket mock** that simulates all stages
- **Handle connection drops** gracefully with polling fallback
- **Test with network interruptions**

**Risk Level**: Medium - Adds complexity but well-understood pattern

---

### Challenge 5: Code Upload Format

**Problem**: How does Backend.im expect code to be uploaded?
- Single file vs. multiple files
- Base64 encoding vs. JSON vs. tar.gz
- Project structure requirements
- Metadata (language, dependencies, etc.)

**Mitigation**:
- **Test with different formats** in mock API
- **Support multiple upload formats** initially
- **Make CLI flexible** - allow configurable upload method
- **Document format assumptions** clearly

**Risk Level**: Low - Easy to adjust once discovered

---

### Challenge 6: Error Handling & Edge Cases

**Problem**: Real API will have errors mock API doesn't cover:
- Rate limiting
- Network timeouts
- Invalid code format
- Deployment failures
- Resource limits

**Mitigation**:
- **Create comprehensive error scenarios** in mock API
- **Test all HTTP status codes** (400, 401, 403, 500, etc.)
- **Implement retry logic** with exponential backoff
- **User-friendly error messages** with actionable guidance
- **Logging and debugging** capabilities

**Risk Level**: Medium - Requires thorough testing

---

### Challenge 7: Testing Without Real Platform

**Problem**: Mock API might not reflect real behavior:
- Performance characteristics
- Actual deployment times
- Real error scenarios
- Integration edge cases

**Mitigation**:
- **Create integration test suite** that runs against both mock and real API
- **Use contract testing** (Pact, etc.) to validate compatibility
- **Early integration testing** - switch to real API as soon as possible
- **Document differences** between mock and real API

**Risk Level**: Low - Can be mitigated with early integration

---

### Challenge 8: Editor Integration (Optional)

**Problem**: Optional editor integration needs to work across platforms:
- Editor command availability (not always in PATH)
- Different editor installation paths
- Windows vs. Linux vs. macOS differences
- Headless environments (no GUI)

**Mitigation**:
- **Make editor integration optional** - CLI works without it
- **Support multiple editors** via `--editor` flag (vscode, vim, nano, etc.)
- **Graceful degradation** - if editor not found, user can open manually
- **Clear error messages** with installation instructions
- **Users can always edit files manually** - no hard dependency

**Risk Level**: Low - Optional feature, easy to skip if problematic

---

## 📊 Overall Risk Assessment

| Challenge | Risk Level | Impact | Mitigation Effort |
|----------|------------|--------|-------------------|
| Incomplete API Discovery | Medium | 2-3 days delay | Low |
| API Contract Drift | Medium | Breaking changes | Medium |
| Google OAuth Authentication | **Low-Medium** | **3-5 days** | **Low** ✅ |
| WebSocket Integration | Medium | 1 week | Medium |
| Code Upload Format | Low | Few days | Low |
| Error Handling | Medium | Ongoing | Medium |
| Testing Limitations | Low | Minor | Low |
| Editor Integration (Optional) | Low | Few days | Low |

**Total Risk Adjustment**: +1-2 weeks to timeline (from 8 to 9-10 weeks)

**Key Improvement**: Google-only OAuth significantly reduces authentication complexity:
- ✅ Well-documented patterns (gcloud auth login style)
- ✅ Official Go libraries available (`golang.org/x/oauth2`)
- ✅ Established best practices for CLI apps
- ✅ No need for multiple provider support

**Recommendation**: The mocking approach is **highly viable** - challenges are manageable and well-understood. Google OAuth simplifies the authentication challenge substantially.

---

## 💰 Cost Analysis

### Additional Infrastructure Costs (NEW - CLI Tool Only)

| Service | Provider | Cost | Purpose |
|---------|----------|------|---------|
| **Go Binary Distribution** | GitHub Releases | Free | CLI tool distribution |
| **Documentation** | GitHub Pages | Free | Documentation hosting |
| **CI/CD** | GitHub Actions | Free | Automated testing and releases |

**Total Additional Cost**: $0/month

### Existing Infrastructure Costs (EXISTING - Backend.im Platform)

| Service | Status | Cost | Purpose |
|---------|--------|------|---------|
| **EC2 Frontend** | Existing | Already paid | Web UI |
| **EC2 Backend** | Existing | Already paid | API services |
| **Kubernetes Cluster** | Existing | Already paid | Container orchestration |
| **Gitea Server** | Existing | Already paid | Git hosting |
| **Database** | Existing | Already paid | User and deployment data |
| **Load Balancer** | Existing | Already paid | Traffic routing |

**Total Existing Cost**: Already covered by Backend.im platform

### Development Costs (NEW - CLI Tool Development)

| Service | Cost | Purpose |
|---------|------|---------|
| **Development Time** | One-time | CLI tool development (4-6 weeks) |
| **Maintenance** | Minimal | Bug fixes and updates |
| **Claude API** | Pay-per-use | AI code generation (existing usage) |

**Total Development Cost**: Minimal (one-time development effort)

---

## 🔒 Security Architecture

### Authentication & Authorization (EXISTING - Backend.im Platform)

```
Developer → Go CLI → OAuth 2.0 → Google → JWT Token → Backend.im API
```

**Security Features (Already Implemented):**
- OAuth 2.0 authentication with Google
- JWT token-based API authentication
- User session management
- Role-based access control

**CLI Implementation**:
- Uses Google OAuth 2.0 installed app flow (similar to `gcloud auth login`)
- Leverages `golang.org/x/oauth2` library
- Secure token storage in local config
- Automatic token refresh handling

### Network Security (EXISTING - Backend.im Platform)

**Security Features (Already Implemented):**
- TLS encryption for all API communications
- Network policies for Kubernetes cluster
- Firewall rules for EC2 instances
- Secrets management in Kubernetes

### Data Protection (EXISTING - Backend.im Platform)

**Security Features (Already Implemented):**
- Encryption at rest for database and storage
- Automated backup strategy
- Comprehensive access logging
- Audit trails for all operations

### CLI Tool Security (NEW - Minimal Requirements)

**Security Features (New):**
- Secure token storage in local configuration
- HTTPS-only communication with Backend.im API
- Input validation for code and prompts
- No sensitive data stored in CLI tool
- Code signing for binary distribution

---

## 📊 Monitoring & Observability

### Existing Monitoring (EXISTING - Backend.im Platform)

**Monitoring Features (Already Implemented):**
- Prometheus metrics collection for all services
- Grafana dashboards for visualization
- Application metrics (response times, error rates)
- Infrastructure metrics (CPU, memory, disk usage)
- Business metrics (deployment frequency, success rates)
- Security metrics (failed authentication attempts)

### CLI Tool Monitoring (NEW - Minimal Requirements)

**Monitoring Features (New):**
- CLI usage metrics (deployments per user)
- Error tracking for CLI operations
- Performance metrics for API calls
- User activity logging

**Implementation:**
```go
// Simple metrics collection in CLI tool
type MetricsCollector struct {
    apiClient *api.Client
}

func (m *MetricsCollector) TrackDeployment(userID string, success bool, duration time.Duration) error {
    return m.apiClient.Post("/metrics/deployment", map[string]interface{}{
        "user_id":   userID,
        "success":   success,
        "duration":  duration.Milliseconds(),
        "timestamp": time.Now().Unix(),
    })
}
```

---

## 🚀 Implementation Roadmap

### Phase 0: API Discovery & Mocking (Week 1)
- [ ] Reverse-engineer Backend.im web UI API calls
- [ ] Document all API contracts (deploy, status, auth)
- [ ] Create mock API server for testing
- [ ] Validate mock API with Web UI behavior
- [ ] **Deliverable**: API_CONTRACTS.md and working mock API

**Benefits**: Enables parallel development and contract-driven approach

### Phase 1: CLI Tool Development (Week 2-3)
- [ ] Create Go CLI tool structure with Cobra
- [ ] Implement authentication with mock API (OAuth flow)
- [ ] Add code generation API client (calls Backend.im `/api/generate`)
- [ ] Implement file download functionality
- [ ] Add deployment API client (calls Backend.im `/api/deploy`)
- [ ] Build against mock API

**Progress**: CLI fully functional with mock API

### Phase 2: File Management & Optional Editor (Week 4)
- [ ] Implement local project file management
- [ ] Add optional editor integration (`--editor` flag)
- [ ] Create deployment command with file upload
- [ ] Add real-time status updates (WebSocket mock)
- [ ] **Deliverable**: Working CLI with mock backend

**Progress**: End-to-end CLI workflow functional locally

### Phase 3: Integration & Testing (Week 5-6)
- [ ] Replace mock API with real Backend.im API
- [ ] End-to-end testing with Backend.im platform
- [ ] Fix any contract mismatches
- [ ] Error handling and edge cases
- [ ] User acceptance testing
- [ ] **Deliverable**: Production-ready CLI

**Progress**: CLI integrated with live Backend.im platform

### Phase 4: Launch & Support (Week 7-8)
- [ ] Documentation and examples
- [ ] Binary distribution setup (GitHub Releases)
- [ ] Public release of CLI tool
- [ ] User onboarding and support
- [ ] Monitor usage and performance
- [ ] Gather feedback for improvements

**Deliverable**: Public release with documentation and support

**Total Timeline: 8 weeks** (With parallel development capability)

---

## 🎯 Success Metrics

### Technical KPIs (NEW - CLI Tool)

| Metric | Target | Measurement |
|--------|--------|-------------|
| **CLI Installation Time** | < 2 minutes | Binary download to working CLI |
| **Deployment Time** | < 2 minutes | Code generation to live URL |
| **CLI Response Time** | < 5 seconds | Command execution time |
| **Error Rate** | < 1% | Failed CLI operations |

### User Experience KPIs (NEW - CLI Tool)

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Learning Curve** | Minimal | Time to first deployment |
| **Developer Satisfaction** | High | User feedback scores |
| **Adoption Rate** | Growing | Active CLI users per month |
| **Command Usage** | High | Frequency of CLI commands |

### Existing Platform KPIs (EXISTING - Backend.im Platform)

| Metric | Status | Measurement |
|--------|--------|-------------|
| **Platform Uptime** | 99.9% | Service availability (already achieved) |
| **API Response Time** | < 1 second | API endpoint performance (already achieved) |
| **Deployment Success Rate** | > 99% | Successful deployments (already achieved) |

---

## 🔄 Future Enhancements

### CLI Tool Enhancements (NEW - Potential Features)

- **Multi-Environment Support**: Deploy to dev, staging, production
- **Custom Domains**: User-provided domains via CLI
- **Database Migrations**: Automated schema updates via CLI
- **Local Development**: Hot reload and local testing
- **Code Templates**: Pre-built templates for common patterns

### Integration Opportunities (NEW - CLI Tool)

- **CI/CD Integration**: GitHub Actions, GitLab CI
- **IDE Integration**: VS Code, IntelliJ plugins
- **Monitoring Integration**: CLI-based monitoring commands
- **Security Scanning**: CLI-based security checks

### Existing Platform Enhancements (EXISTING - Backend.im Platform)

**Note**: These features would be added to the existing Backend.im platform, not the CLI tool:
- **Auto-scaling**: Based on traffic patterns
- **Blue-Green Deployments**: Zero-downtime updates
- **Advanced Monitoring**: Enhanced observability features

---

## 📝 Conclusion

This research presents a minimal, cost-effective solution for CLI-based deployment to Backend.im. The proposed approach leverages the existing Backend.im platform infrastructure while adding a simple Go-based CLI tool for seamless developer experience.

**Key Advantages:**
- ✅ **Cost-Effective**: $0 additional infrastructure cost
- ✅ **Minimal Development**: Only CLI tool required
- ✅ **Leverages Existing**: Uses all existing Backend.im services
- ✅ **Developer-Friendly**: Simple binary installation and go
- ✅ **Secure**: Reuses existing authentication and security
- ✅ **Fast Implementation**: 6-8 week development timeline

**What's New:**
- Go-based CLI tool for developers
- Local file management for generated code
- Optional editor integration for code editing
- CLI-based deployment workflow that leverages Backend.im's existing generation

**What's Existing:**
- Backend.im platform (EC2, K8s, Gitea, API, WebSocket)
- All security, monitoring, and infrastructure
- User authentication and management
- Deployment orchestration

The solution provides a simple, effective way to enable CLI-based deployment to Backend.im without requiring any changes to the existing platform infrastructure.

---

**Research Completed By**: [Your Name]  
**Date**: [Current Date]  
**Version**: 1.0
