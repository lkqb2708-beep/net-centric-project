# MangaHub Project Brief and Prompt Pack

## 1) Source documents reviewed
- MangaHub CLI Application User Manual
- MangaHub Use Case Specification (Revised Scope)
- MangaHub Manga & Comic Tracking System Term Project Description (Revised Scope)

## 2) Key details extracted from the files

### Project identity
- Project name: **MangaHub**
- Domain: **Manga and comic tracking system**
- Course context: **Net-centric programming / network programming**
- Required language: **Go**
- Team size: **2 students**
- Timeline: **10–12 weeks**
- Educational goal: demonstrate practical network programming with **HTTP, TCP, UDP, gRPC, and WebSocket**.

### Core purpose
MangaHub is a manga tracking platform with:
- user registration and authentication
- manga search and details
- personal reading library
- reading progress tracking
- real-time sync
- chapter notifications
- live chat
- analytics/statistics
- export/import and backup features
- admin/server health and troubleshooting tooling

### System requirements and platform notes
- Go 1.19+ in the manual
- UTF-8 terminal support
- Network connectivity for sync features
- Linux / macOS / Windows support
- CLI-first architecture in the manual, but the spec also defines web/API-style endpoints and a simplified UI scope
- Storage in the spec is simplified around a local database plus JSON data files for manga seed content

### Data / content requirements
The project spec asks for a manageable manga dataset:
- at least **30–40 manga series**
- basic metadata: title, author, genres, status, chapter count, description, cover URL
- data sources can include manual entry, limited API integration, and lightweight scraping for practice
- simplified user progress structure with reading lists and per-title chapter tracking

### Main application areas
#### Authentication
- Register
- Login
- Logout
- Check auth status
- Change password
- JWT-based authorization in the spec
- password hashing with bcrypt in the use cases

#### Manga discovery and management
- Search manga
- View manga details
- List manga
- Advanced filtering by genre, status, author, year, chapter count, sort order
- Add manga to personal library
- Update/remove library entries
- Track reading status: reading, completed, plan-to-read, on-hold, dropped

#### Progress tracking
- Update current chapter and volume
- Prevent invalid chapter progress
- Support history and sync
- Broadcast progress updates to connected clients

#### Network protocols
The revised scope explicitly requires these five protocols:
- **HTTP REST API**
- **TCP progress sync**
- **UDP notifications**
- **WebSocket chat**
- **gRPC internal service**

#### Chat system
- real-time general chat
- manga-specific chat rooms
- join/leave behavior
- list online users
- private messages
- chat history
- user presence/status

#### Server management / ops
- start/stop services
- server status and health checks
- logs and debug modes
- ping/connectivity tests
- recovery from connection failure and database issues
- graceful degradation when services are unavailable

#### Statistics, export, backup
- reading statistics dashboard
- reading trends and goals
- export library/progress/all data
- backup and restore
- database check/repair/optimize/stats

### Simplified backend schema in the spec
The spec defines a small relational model:
- users
- manga
- user_progress

The manual and use cases also add:
- reviews and ratings
- friends and activity feed
- cached popular manga data
- notification preferences
- multiple reading lists
- advanced search/filtering
- optional Redis caching

### Performance, reliability, and security targets
- support around **50–100 concurrent users**
- search response targets around **500 ms**
- progress sync broadcast within **1 second**
- chat delivery around **100 ms**
- proper concurrency handling
- JWT validation
- input sanitization
- graceful shutdown
- connection cleanup
- error recovery for TCP/WebSocket/database issues

### Grading / project emphasis
The spec places the highest value on:
- working implementation of all five protocols
- integration of services
- database persistence
- code quality and testing
- documentation and live demo readiness

### Academic integrity note
The project description explicitly says AI tools may be used for brainstorming, refinement, and summarization, but final implementation and submission must be the students’ own work. Any AI usage must be acknowledged in the report.

---

## 3) Recommended implementation stance for a Go build
Because the project language is Go, the safest architecture is:
- Go backend services
- PostgreSQL database
- pgAdmin for database administration
- Go-native database access and migration tooling
- avoid putting a Node/TypeScript ORM into the Go backend layer

## 4) Gemini / Claude prompt

### Copy/paste prompt
You are a senior Go engineer and product architect. Build a complete, production-style academic project named **MangaHub**, a manga and comic tracking system for a network programming course.

Requirements:
- Language: **Go**
- Project scope: a full working prototype with a polished UI, complete animations/micro-interactions, auth screens, logout flow, database persistence, and all requested services.
- Must support these protocols:
  - HTTP REST API
  - TCP progress synchronization
  - UDP notification broadcasting
  - gRPC internal service
  - WebSocket real-time chat
- The app should include:
  - landing page / main menu
  - login screen
  - registration screen
  - logout flow
  - manga search
  - manga detail page
  - library page
  - reading progress update flow
  - statistics dashboard
  - notifications area
  - chat area
  - settings / profile area
  - server health / status view
- Use a complete, responsive UI with:
  - modern layout
  - clear navigation
  - animated page transitions
  - loading states
  - empty states
  - error states
  - toast/alert feedback
  - accessible controls
- Use a real database layer:
  - PostgreSQL as the database
  - include migration/setup scripts
  - connect the database to a visual admin tool such as pgAdmin
  - define tables for users, manga, user_progress, and any additional tables needed for reviews, chat, friends, notifications, and stats
- Use secure auth:
  - register
  - login
  - logout
  - JWT session handling
  - password hashing
  - protected routes
- Implement the full feature set from the specification:
  - manga search and filtering
  - manga details
  - library management
  - progress tracking
  - progress history
  - TCP broadcast sync
  - UDP release notifications
  - WebSocket chat
  - gRPC internal calls
  - server status and diagnostics
  - export/backup utilities
- Build the codebase in a clean modular structure:
  - cmd/
  - internal/
  - pkg/
  - proto/
  - data/
  - docs/
- Include a frontend and backend that feel like one cohesive product.
- Prefer realistic implementation over fake mock-only behavior.
- Provide:
  - base project scaffold
  - file tree
  - database schema
  - API routes
  - service architecture
  - UI components
  - animation behavior
  - step-by-step setup instructions
  - run instructions
  - environment variables
  - test strategy
  - demo checklist
- If something is too large for one pass, still generate a strong base plate and the complete architecture so the project is immediately extendable.
- Make the result suitable for a class demo and for further student refinement.

### Positive prompt
- modern manga tracker
- polished responsive UI
- animated transitions
- clean auth flow
- real database persistence
- Go backend
- multi-protocol architecture
- practical academic scope
- modular codebase
- demo-ready
- clear error handling
- secure authentication
- real-time chat and sync
- production-style structure

### Negative prompt
- no placeholder-only UI
- no broken routing
- no hardcoded fake auth
- no in-memory-only database
- no unfinished login/logout flow
- no missing navigation
- no static screenshots instead of a working app
- no single-file monolith
- no mixing random tech stacks that conflict with Go
- no Prisma inside the Go application layer
- no incomplete protocol stubs without usable behavior
- no unreadable code or toy-level structure
- no skipping server status, error handling, or cleanup

## 5) Suggested database choice
Use **PostgreSQL** for the application database and manage it with **pgAdmin**. In a Go backend, use Go-native database access and migrations rather than trying to force a Node/TypeScript ORM into the stack.

## 6) Implementation checklist
- Project scaffold
- Auth
- UI shell and navigation
- Manga catalog
- Library management
- Progress tracking
- TCP sync server
- UDP notification server
- WebSocket chat server
- gRPC service
- stats dashboard
- backup/export tools
- logs and health checks
- documentation
- demo script

## 7) Suggested report contents
If you hand this to a teammate or model, the final report should explain:
- what was implemented
- what remains optional
- what tech stack was chosen
- how to run the project
- how the protocols interact
- how the database is structured
- how login/logout works
- how the UI is organized
- how to demo the system
