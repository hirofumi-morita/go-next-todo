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
│   ├── handlers/     # API route handlers
│   ├── middleware/   # Authentication middleware
│   ├── models/       # Database models
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

2. **TODO Management (CRUD)**
   - Create, Read, Update, Delete TODOs
   - Mark TODOs as complete/incomplete
   - Each user sees only their own TODOs

3. **Admin User Management**
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
- `POST /api/todos` - Create new TODO
- `PUT /api/todos/:id` - Update TODO
- `DELETE /api/todos/:id` - Delete TODO

### Admin Routes (require admin role)
- `GET /api/admin/users` - Get all users
- `GET /api/admin/users/:id` - Get specific user with TODOs
- `DELETE /api/admin/users/:id` - Delete user
- `PATCH /api/admin/users/:id` - Update user admin status

## Environment Variables
- `DATABASE_URL` - PostgreSQL connection string
- `SESSION_SECRET` - JWT signing secret

## Running the Application
The application runs two services:
1. Go backend API on port 8080
2. Next.js frontend on port 5000

## Recent Changes
- Initial project setup (November 28, 2025)
