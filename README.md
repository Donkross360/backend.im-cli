# Backend.im CLI

A command-line tool for generating, editing, and deploying backend code to the Backend.im platform.

## Features

- 🔐 **Authentication** - Login with Backend.im (Google OAuth)
- 🚀 **Code Generation** - Generate FastAPI code from natural language prompts
- ✏️ **Code Editing** - Edit generated code locally with your favorite editor
- 💾 **Commit Changes** - Save your edits back to Backend.im/Gitea
- 📦 **Deploy** - Deploy your code with real-time WebSocket progress updates

## Quick Start

### Option 1: Local Installation (Recommended)

**Install the CLI:**
```bash
./install.sh
```

**Configure for mock API:**
```bash
export BACKEND_IM_API_URL=http://localhost:8080
echo 'export BACKEND_IM_API_URL=http://localhost:8080' >> ~/.bashrc
```

**Start mock API:**
```bash
docker-compose up -d mock-api
```

**Run your first deployment:**
```bash
backend-im login
backend-im generate "Create a REST API" --project my-api --output ./my-api
backend-im commit my-api --dir ./my-api
backend-im deploy my-api --dir ./my-api
```

### Option 2: Docker (Development)

**Start services:**
```bash
make docker-up          # Start mock API
make docker-build       # Build CLI
```

**Run commands:**
```bash
docker-compose run --rm cli login
docker-compose run --rm cli generate "Create a REST API" --project my-api
docker-compose run --rm cli commit my-api
docker-compose run --rm cli deploy my-api
```

## Installation

### Local Binary Installation

The easiest way is using the install script (no Go installation needed):

```bash
cd backenb.im-cli
./install.sh
```

This will:
- Build the binary using Docker
- Install to `/usr/local/bin/backend-im`
- Create config directory at `~/.backend-im`

**Custom install location:**
```bash
INSTALL_DIR=~/bin ./install.sh
```

See [LOCAL_INSTALL.md](./LOCAL_INSTALL.md) for detailed instructions.

### Docker Installation

Build and run everything in Docker:

```bash
# Build images
make docker-build

# Start mock API
make docker-up

# Run CLI commands
docker-compose run --rm cli [command]
```

## Commands

### `login` - Authenticate with Backend.im

First-time setup - authenticate with Backend.im.

```bash
backend-im login
```

**Alias:** `auth` (for backward compatibility)

---

### `generate` - Generate Code from Prompt

Generate FastAPI code from a natural language prompt. Code is automatically committed to Backend.im/Gitea.

```bash
backend-im generate "Create a REST API for user management" --project my-api --output ./my-api
```

**Options:**
- `--project, -p` - Project ID (required)
- `--output, -o` - Output directory (default: auto-generated name)
- `--editor, -e` - Auto-open in editor (default: uses $EDITOR if set)

**Example:**
```bash
backend-im generate "Create a FastAPI app with SQLAlchemy" \
  --project my-api \
  --output ./my-api \
  --editor code
```

---

### `edit` - Open Project in Editor

Open a project directory in your default editor to view and edit code.

```bash
backend-im edit [project-directory]
```

**Options:**
- `--dir, -d` - Project directory (default: current directory)
- `--editor, -e` - Editor to use (default: $EDITOR or auto-detect)

**Example:**
```bash
backend-im edit ./my-api
backend-im edit ./my-api --editor vim
```

---

### `commit` - Commit Local Changes

Save your local file changes to Backend.im/Gitea.

```bash
backend-im commit [project-id] --dir ./my-api --message "My changes"
```

**Options:**
- `[project-id]` - Project ID (positional argument or `--project, -p`)
- `--dir, -d` - Project directory (default: current directory)
- `--message, -m` - Commit message (default: "Update code from CLI")

**Examples:**
```bash
# Using positional argument
backend-im commit my-api --dir ./my-api

# Using flag
backend-im commit --project my-api --dir ./my-api --message "Fixed bug"

# From project directory
cd ./my-api
backend-im commit my-api
```

---

### `deploy` - Deploy to Backend.im

Deploy your code to Backend.im. Automatically watches deployment progress via WebSocket.

```bash
backend-im deploy [project-id] --dir ./my-api
```

**Options:**
- `[project-id]` - Project ID (positional argument or `--project, -p`)
- `--dir, -d` - Project directory (default: current directory)
- `--watch, -w` - Watch deployment progress (default: true, use `--watch=false` to disable)

**Examples:**
```bash
# Using positional argument (recommended)
backend-im deploy my-api --dir ./my-api

# Without watching (faster, but no progress updates)
backend-im deploy my-api --dir ./my-api --watch=false

# From project directory
cd ./my-api
backend-im deploy my-api
```

**Deployment Progress:**
The deploy command streams real-time updates showing:
- Status changes (committing → creating_namespace → deploying → complete)
- Logs from the deployment process
- Final deployment URL when complete

---

## Complete Workflow Example

```bash
# 1. Login (first time only)
backend-im login

# 2. Generate code from prompt
backend-im generate "Create a user management API with authentication" \
  --project user-api \
  --output ./user-api

# 3. Edit the generated code (optional)
backend-im edit ./user-api
# Make your changes...

# 4. Commit your changes
backend-im commit user-api --dir ./user-api --message "Added user validation"

# 5. Deploy
backend-im deploy user-api --dir ./user-api
# Watch the real-time progress and get the deployment URL!
```

## Configuration

### Environment Variables

- `BACKEND_IM_API_URL` - API endpoint URL (default: `http://localhost:8080`)

**For local mock API:**
```bash
export BACKEND_IM_API_URL=http://localhost:8080
```

**For production:**
```bash
export BACKEND_IM_API_URL=https://api.backend.im
```

### Config Directory

The CLI stores configuration in `~/.backend-im/`:
- `token.json` - Authentication token (automatically managed)

## Project Structure

```
backenb.im-cli/
├── cli/                      # CLI tool source code
│   ├── cmd/backend-im/      # Main entry point
│   └── internal/
│       ├── api/             # API client
│       ├── auth/             # Authentication
│       ├── commands/         # CLI commands (login, generate, commit, deploy, edit)
│       ├── editor/           # Editor integration
│       └── files/            # File operations
├── mock-api/                 # Mock Backend.im API server (for testing)
├── docker-compose.yml        # Docker setup
├── install.sh                # Local installation script
├── Makefile                  # Build commands
└── README.md                 # This file
```

## Development

### Prerequisites

- Docker (for building without Go)
- Or Go 1.21+ (for local development)

### Build Commands

```bash
# Build CLI binary using Docker
make docker-build

# Start mock API
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs

# Run CLI commands via Docker
make cli-help
make cli-auth
make cli-generate
```

### Testing

**Test mock API endpoints:**
```bash
# Start mock API
make docker-up

# Test code generation
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Create a REST API"}'

# Test deployment
curl -X POST http://localhost:8080/api/deploy \
  -H "Content-Type: application/json" \
  -d '{"files": {"main.py": "print(\"hello\")"}, "projectId": "test-project"}'
```

## Troubleshooting

### Permission Issues

**Problem:** "Permission denied" when saving token
```bash
sudo chown $USER:$USER ~/.backend-im
chmod 700 ~/.backend-im
```

### Connection Issues

**Problem:** Can't connect to mock API
```bash
# Check if mock API is running
docker-compose ps mock-api

# Verify API URL is set
echo $BACKEND_IM_API_URL

# Set it if missing
export BACKEND_IM_API_URL=http://localhost:8080
```

### Binary Not Found

**Problem:** `backend-im: command not found`
```bash
# Check if installed
which backend-im

# Verify installation
ls -l /usr/local/bin/backend-im

# Reinstall if needed
./install.sh
```

## Status

### ✅ Completed
- Authentication (login/auth commands)
- Code generation with prompt
- Local code editing
- Commit changes to Backend.im
- Deploy with WebSocket streaming
- Real-time deployment progress
- Local binary installation
- Docker setup

### 🚧 Future Enhancements
- Real Google OAuth integration (currently mock)
- Token refresh handling
- Project listing and management
- Deployment history
- Environment variable management

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
