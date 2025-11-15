# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Agent Platform (agent-opus) is an intelligent agent management and orchestration system built with a **dual-architecture approach**:
- **Backend**: gRPC services with HTTP Gateway (grpc-gateway) for REST API compatibility
- **Frontend**: React 18 + Vite + TypeScript with Ant Design UI

The platform provides Agent management, conversation handling, tool integration, and knowledge base capabilities.

## Architecture

### Backend: gRPC + HTTP Gateway Pattern

The backend uses **gRPC as the primary service layer** with an HTTP Gateway for REST compatibility:

1. **Proto Definitions** (`proto/*.proto`): Define all service contracts
2. **gRPC Services** (`backend/internal/grpc/`): Core business logic implementation
3. **HTTP Gateway** (`backend/cmd/server/gateway.go`): Automatic REST API generation via grpc-gateway
4. **Dual Protocol Support**:
   - gRPC on port 9000 (internal services)
   - HTTP/REST on port 8000 (external API)

**Key Implementation Details:**
- Both gRPC and HTTP servers run concurrently from `cmd/server/main.go`
- HTTP Gateway automatically converts REST calls to gRPC calls
- Custom marshaler (`cmd/server/marshaler.go`) ensures consistent response format: `{code, message, data, timestamp, request_id}`
- CORS middleware in `gateway.go` handles cross-origin requests

### Data Layer

- **ORM**: ent (schema-based code generation)
- **Model Schemas**: `backend/internal/model/schema/` - Define database entities
- **Databases**: MySQL (primary), MongoDB (documents), Redis (cache), Milvus (vectors)

### Frontend Architecture

- **State Management**: Redux Toolkit (`frontend/src/store/`)
- **Routing**: React Router v6
- **API Layer**: `frontend/src/services/` - Centralized API calls
- **Component Structure**:
  - `frontend/src/pages/` - Feature pages (Agent, Conversation, Tool, KnowledgeBase, Dashboard)
  - `frontend/src/components/Layout/` - Layout components

## Common Commands

### Backend Development

```bash
# Start backend locally (no database required - uses mock data)
cd backend
go run cmd/server/main.go

# Build backend binary
cd backend
go build -o agent-platform cmd/server/*.go

# Or use Makefile
make backend

# Run backend tests
cd backend
go test ./...

# Format backend code
cd backend
go fmt ./...
```

### Frontend Development

```bash
# Install dependencies
cd frontend
npm install

# Start dev server (http://localhost:5173)
npm run dev

# Build for production
npm run build

# Or use Makefile
make frontend

# Lint frontend code
cd frontend
npm run lint
```

### Proto Generation

When modifying `.proto` files, regenerate code:

```bash
# Install proto tools first (one-time)
make install-proto-tools

# Generate proto code for both backend and frontend
make proto

# Or separately
make proto-backend    # Go code generation
make proto-frontend   # TypeScript code generation

# Clean generated proto files
make proto-clean
```

**Important**: After proto changes:
1. Run `make proto-backend` to regenerate Go code in `backend/gen/go/`
2. Restart the backend server
3. Run `make proto-frontend` for frontend updates (requires protoc-gen-grpc-web)

### Docker Operations

```bash
# Start all services
make start
# or
docker-compose up -d

# View logs
make logs
# or
docker-compose logs -f backend

# Stop services
make stop

# Clean up containers and volumes
make clean
```

### Testing

```bash
# Test all
make test

# Test backend only
cd backend && go test ./...

# Test API endpoints (requires backend running)
./test-api.sh

# Quick API check
curl http://localhost:8000/health
curl http://localhost:8000/api/v1/ping
```

## Development Workflow

### Adding a New Service Feature

1. **Define Proto Contract**: Add service definition in `proto/` directory
2. **Generate Code**: Run `make proto-backend`
3. **Implement gRPC Service**: Create service in `backend/internal/grpc/`
4. **Register Service**: Add registration in `cmd/server/main.go` and `gateway.go`
5. **HTTP Gateway Auto-generates REST Endpoints**: No manual REST handler needed

Example service structure:
- Proto: `proto/my_service.proto`
- Implementation: `backend/internal/grpc/my_service.go`
- Auto-generated: HTTP endpoints via grpc-gateway

### Adding Database Models

1. Define schema in `backend/internal/model/schema/`
2. Run `go generate ./...` in backend directory to generate ent code
3. Access via repositories in `backend/internal/repository/`

### Adding Frontend Pages

1. Create page component in `frontend/src/pages/`
2. Add route in `frontend/src/App.tsx`
3. Create API service in `frontend/src/services/`
4. Define TypeScript types in `frontend/src/types/`

## Configuration

### Backend Environment Variables

Key configurations in `backend/.env`:

```env
# Server ports
GRPC_PORT=9000          # gRPC server port
HTTP_PORT=8000          # HTTP Gateway port
SERVER_MODE=debug       # debug or release

# Databases
MYSQL_HOST=localhost
MYSQL_PORT=3306
REDIS_HOST=localhost
REDIS_PORT=6379
MONGODB_URI=mongodb://localhost:27017

# AI Integration
OPENAI_API_KEY=your-key
EMBEDDING_MODEL=text-embedding-ada-002

# Security
JWT_SECRET=your-secret-key
CORS_ORIGINS=http://localhost:5173
```

### Port Conflicts

If local services conflict with Docker ports:
- Modify `docker-compose.yml` port mappings (e.g., `3307:3306` for MySQL)
- Update corresponding env vars in `backend/.env`
- Or run databases only: `docker-compose up -d mysql mongodb redis`

## Project Structure Highlights

```
backend/
├── cmd/server/
│   ├── main.go          # gRPC server + HTTP Gateway startup
│   ├── gateway.go       # HTTP Gateway setup & CORS
│   └── marshaler.go     # Custom response marshaler
├── internal/
│   ├── grpc/            # gRPC service implementations
│   ├── config/          # Configuration loading
│   ├── middleware/      # Auth, CORS, logging
│   ├── model/schema/    # ent schemas
│   ├── repository/      # Data access layer
│   ├── response/        # Response formatting
│   └── service/         # Business logic (if needed beyond gRPC)
└── gen/go/              # Generated protobuf/gRPC code

proto/                   # Proto definitions for all services
├── agent.proto
├── conversation.proto
├── tool.proto
├── knowledge_base.proto
└── common.proto

frontend/
├── src/
│   ├── pages/           # Feature pages
│   ├── components/      # Reusable components
│   ├── services/        # API client layer
│   ├── store/           # Redux state management
│   └── types/           # TypeScript definitions
└── vite.config.ts       # Vite config with proxy to backend
```

## Current Implementation Status

**MVP Complete**:
- ✅ Backend using mock data (no database connection required)
- ✅ Full CRUD APIs for Agents, Conversations, Tools, Knowledge Bases
- ✅ Frontend UI with all main pages
- ✅ Docker containerization
- ✅ API testing scripts

**To Implement**:
- Database persistence (MySQL + ent ORM integration)
- Real AI model integration (OpenAI/Anthropic API calls)
- Vector database (Milvus) for knowledge base
- Multi-Agent workflow orchestration
- Authentication & authorization

## API Response Format

All HTTP Gateway responses follow this format:

```json
{
  "code": 0,
  "message": "success",
  "data": { /* actual response data */ },
  "timestamp": 1634567890,
  "request_id": "uuid"
}
```

This is enforced by the custom marshaler in `cmd/server/marshaler.go`.

## Key Dependencies

**Backend**:
- `google.golang.org/grpc` - gRPC framework
- `grpc-ecosystem/grpc-gateway` - HTTP Gateway
- `entgo.io/ent` - ORM
- `gin-gonic/gin` - HTTP utilities (limited use)
- `go.uber.org/zap` - Structured logging

**Frontend**:
- `react` 18 + `react-router-dom` v6
- `antd` 5 (Ant Design)
- `@reduxjs/toolkit` - State management
- `axios` - HTTP client
- `reactflow` - Workflow visualization

## Testing Notes

- Backend currently uses mock data, allowing API testing without database setup
- Start backend: `cd backend && go run cmd/server/main.go`
- Test APIs: `./test-api.sh` or manual curl commands
- Frontend connects to backend via Vite proxy (configured in `vite.config.ts`)
