# Quick Start Guide

Get up and running with Backend.im CLI in 5 minutes!

## Prerequisites

- Docker (for mock API and building binary)
- Linux/macOS system

## Step 1: Install CLI

```bash
cd backenb.im-cli
./install.sh
```

## Step 2: Configure

```bash
# Set API URL for mock API
export BACKEND_IM_API_URL=http://localhost:8080

# Make it permanent
echo 'export BACKEND_IM_API_URL=http://localhost:8080' >> ~/.bashrc
source ~/.bashrc
```

## Step 3: Start Mock API

```bash
docker-compose up -d mock-api
```

## Step 4: Login

```bash
backend-im login
```

## Step 5: Your First Deployment

```bash
# Generate code
backend-im generate "Create a simple hello world API" \
  --project hello-api \
  --output ./hello-api

# Commit (optional - code is already committed)
backend-im commit hello-api --dir ./hello-api

# Deploy
backend-im deploy hello-api --dir ./hello-api
```

That's it! 🎉

## Common Commands

```bash
# Login
backend-im login

# Generate code
backend-im generate "Your prompt" --project my-api --output ./my-api

# Edit code
backend-im edit ./my-api

# Commit changes
backend-im commit my-api --dir ./my-api --message "My changes"

# Deploy
backend-im deploy my-api --dir ./my-api
```

## Need Help?

- Full documentation: See [README.md](./README.md)
- Local installation details: See [LOCAL_INSTALL.md](./LOCAL_INSTALL.md)
- Troubleshooting: See README.md Troubleshooting section

