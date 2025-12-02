# TODO Management Application

## Overview
A full-stack TODO management application with user authentication and admin user management features.

## Technology Stack
- **Frontend**: Next.js 14+ (App Router) with TypeScript and Tailwind CSS
- **Backend**: Go with Gin framework and GORM
- **Database**: PostgreSQL

## Project Structure
```
/
├── backend/           # Go backend API server
│   ├── config/       # Database configuration
│   ├── graph/        # GraphQL layer (gqlgen)
│   │   ├── schema.graphqls  # GraphQL schema
│   │   ├── resolver.go      # Resolver initialization
│   │   ├── schema.resolvers.go  # Query/Mutation resolvers
│   │   ├── client.go        # GraphQL client for REST handlers
│   │   └── model/           # Generated models
│   ├── handlers/     # REST API route handlers (use GraphQL)
│   ├── middleware/   # Authentication middleware
│   ├── models/       # Database models (GORM)
│   └── main.go       # Entry point
├── frontend/          # Next.js frontend
│   ├── app/          # App Router pages
│   ├── components/   # React components
│   └── lib/          # API utilities and context
└── replit.md         # This file
```

## Features
1. **User Registration & Login**
   - Email and password authentication
   - JWT token-based session management
   - Password hashing with bcrypt
   - First registered user automatically becomes admin

2. **TODO Management (CRUD)**
   - Create, Read, Update, Delete TODOs
   - Mark TODOs as complete/incomplete
   - Assign TODOs to groups
   - Filter TODOs by group
   - Each user sees only their own TODOs

3. **Group Management**
   - Create, Read, Update, Delete groups
   - Custom group colors
   - Group descriptions
   - TODOs can be assigned to groups

4. **Admin User Management**
   - View all users
   - Delete users
   - Grant/revoke admin privileges

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Protected Routes (require JWT token)
- `GET /api/me` - Get current user info

### TODO Routes
- `GET /api/todos` - Get all TODOs for user
- `GET /api/todos/:id` - Get specific TODO
- `POST /api/todos` - Create new TODO (with optional group_id)
- `PUT /api/todos/:id` - Update TODO (including group assignment)
- `DELETE /api/todos/:id` - Delete TODO

### Group Routes
- `GET /api/groups` - Get all groups for user
- `GET /api/groups/:id` - Get specific group
- `POST /api/groups` - Create new group (name, description, color)
- `PUT /api/groups/:id` - Update group
- `DELETE /api/groups/:id` - Delete group (unlinks TODOs from group)

### Admin Routes (require admin role)
- `GET /api/admin/users` - Get all users
- `GET /api/admin/users/:id` - Get specific user with TODOs
- `DELETE /api/admin/users/:id` - Delete user
- `PATCH /api/admin/users/:id` - Update user admin status

## Environment Variables
- `DATABASE_URL` - PostgreSQL connection string
- `SESSION_SECRET` - JWT signing secret (required)

## Running the Application

### On Replit
The application runs two services:
1. Go backend API on port 8080
2. Next.js frontend on port 5000

### Docker (Local Development)
```bash
# Production build
docker-compose up --build

# Development with hot-reload
docker-compose -f docker-compose.dev.yml up --build
```
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## Project Structure (Docker)
```
/
├── docker/
│   ├── backend.Dockerfile      # Production Go build
│   ├── frontend.Dockerfile     # Production Next.js build
│   ├── backend.dev.Dockerfile  # Dev build with Air hot-reload
│   ├── frontend.dev.Dockerfile # Dev build with hot-reload
│   └── README.md              # Docker setup instructions
├── docker-compose.yml          # Production compose
├── docker-compose.dev.yml      # Development compose
└── .env.example               # Environment variable template
```

## Admin Bootstrap
The first user to register automatically receives admin privileges.

## Architecture

### GraphQL Layer
The backend uses gqlgen to provide a GraphQL layer between REST handlers and the database:
- **Schema**: Defines User, Todo, Group types and Query/Mutation operations
- **Resolvers**: Implement all database operations via GORM
- **Client**: Provides a clean interface for REST handlers to call GraphQL operations

```
REST Handler → GraphQL Client → GraphQL Resolver → GORM → PostgreSQL
```

## Recent Changes
- Initial project setup (November 28, 2025)
- Added security improvements: SESSION_SECRET validation, JWT signing method verification
- Added Docker Compose support for local development (November 28, 2025)
- Added GraphQL layer (gqlgen) for database operations (December 1, 2025)
- Updated Docker Compose to support GraphQL code generation (December 1, 2025)
- Added Group feature for organizing TODOs (December 2, 2025)
- Fixed authorization vulnerabilities in group operations (December 2, 2025):
  - CreateTodo/UpdateTodo now verify group ownership before assignment
  - DeleteGroup now verifies ownership before unlinking todos
  - Frontend only sends group_id when explicitly changed
