# Local Installation Guide

## Installation

Run the install script:

```bash
./install.sh
```

This will:
- Build the binary using Docker (no Go installation needed)
- Install to `/usr/local/bin/backend-im`
- Create config directory at `~/.backend-im`

## Configuration

To use with the local mock API, set the environment variable:

```bash
export BACKEND_IM_API_URL=http://localhost:8080
```

**Make it permanent** - Add to your `~/.bashrc` or `~/.zshrc`:

```bash
echo 'export BACKEND_IM_API_URL=http://localhost:8080' >> ~/.bashrc
source ~/.bashrc
```

## Usage

### 1. Start Mock API (if not running)
```bash
cd /path/to/backenb.im-cli
docker-compose up -d mock-api
```

### 2. Login
```bash
backend-im login
```

### 3. Generate Code
```bash
backend-im generate "Create a REST API" --project my-api --output ./my-api
```

### 4. Edit Code (optional)
```bash
backend-im edit my-api
# Or manually edit files
```

### 5. Commit Changes
```bash
backend-im commit my-api --dir ./my-api --message "My changes"
```

### 6. Deploy
```bash
backend-im deploy my-api --dir ./my-api
```

## Shorter Syntax

All commands support positional arguments:

```bash
# Instead of: backend-im deploy --project my-api
backend-im deploy my-api

# Instead of: backend-im commit --project my-api
backend-im commit my-api
```

## Verify Installation

```bash
# Check if installed
which backend-im

# Test commands
backend-im --help
backend-im login
```

## Troubleshooting

**Problem**: "Permission denied" when saving token
**Solution**: Fix config directory permissions:
```bash
sudo chown $USER:$USER ~/.backend-im
chmod 700 ~/.backend-im
```

**Problem**: Can't connect to mock API
**Solution**: Make sure mock API is running and BACKEND_IM_API_URL is set:
```bash
docker-compose ps mock-api
export BACKEND_IM_API_URL=http://localhost:8080
```

**Problem**: Binary not found
**Solution**: Check installation path:
```bash
ls -l /usr/local/bin/backend-im
# If missing, re-run install.sh
```

