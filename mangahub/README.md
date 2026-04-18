# MangaHub 🎌

> A full-stack manga tracking and community platform built in Go.  
> Net-Centric Programming Course Project — demonstrates HTTP, TCP, UDP, WebSocket, and gRPC.

---

## 📋 Features

- ✅ JWT authentication (register, login, logout, profile)
- ✅ Manga catalog with 35+ pre-seeded titles
- ✅ Personal library with reading status and progress tracking
- ✅ Reading history timeline
- ✅ Real-time chat via WebSocket
- ✅ UDP chapter release notifications
- ✅ TCP reading-progress sync broadcast
- ✅ gRPC internal service (MangaService)
- ✅ Statistics dashboard with export
- ✅ Admin panel with health monitoring and logs
- ✅ Polished dark-mode UI with animations

---

## 🏗️ Project Structure

```
mangahub/
├── cmd/server/          Main HTTP server entry point
├── internal/
│   ├── auth/            JWT + bcrypt
│   ├── config/          Environment loading
│   ├── db/              PostgreSQL pool + migrations
│   ├── grpc/            gRPC server implementation
│   ├── handlers/        HTTP request handlers
│   ├── middleware/       Auth, Logger, CORS
│   ├── models/          Domain structs + DTOs
│   ├── realtime/        WebSocket hub
│   ├── repositories/    Database queries
│   ├── router/          Route definitions
│   ├── tcp/             TCP sync server
│   └── udp/             UDP notification server
├── proto/               .proto + pre-generated Go stubs
├── data/
│   ├── migrations/      SQL migration files (001–007)
│   └── seeds/           Optional manga.json seed
├── web/
│   ├── static/          CSS, JS, images
│   └── templates/pages  HTML pages (15 screens)
├── .env                 Development environment variables
├── go.mod               Go module
└── Makefile             Build/run/migrate helpers
```

---

## ⚡ Quick Start

### Prerequisites
- **Go 1.22+** — https://go.dev/dl/
- **PostgreSQL 14+** — https://www.postgresql.org/download/
- **pgAdmin 4** (optional) — https://www.pgadmin.org/

### 1. Clone / open the project
```
cd "C:\Users\Dell\Desktop\NetCen Proj\mangahub"
```

### 2. Create the PostgreSQL database
Open **psql** or **pgAdmin** and run:
```sql
CREATE USER mangahub WITH PASSWORD 'mangahub_password';
CREATE DATABASE mangahub_db OWNER mangahub;
GRANT ALL PRIVILEGES ON DATABASE mangahub_db TO mangahub;
```

### 3. Configure environment
The `.env` file is already created with dev defaults. Edit if your DB credentials differ.

### 4. Download Go dependencies
```powershell
# PowerShell (run in mangahub/ directory)
& "C:\Program Files\Go\bin\go.exe" mod tidy
```

### 5. Run migrations + seed data
```powershell
& "C:\Program Files\Go\bin\go.exe" run ./cmd/server -migrate
& "C:\Program Files\Go\bin\go.exe" run ./cmd/server -seed
```
This creates all tables and inserts 35 manga + admin/demo accounts.

### 6. Start the server
```powershell
& "C:\Program Files\Go\bin\go.exe" run ./cmd/server
```

All four services start together:
| Service | Port | Protocol |
|---------|------|----------|
| HTTP + WebSocket | 8080 | HTTP/WS |
| TCP Sync | 9001 | TCP |
| UDP Notify | 9002 | UDP |
| gRPC | 9003 | gRPC |

### 7. Open the app
Navigate to → **http://localhost:8080**

---

## 🔑 Demo Accounts

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@mangahub.dev | admin1234 |
| User  | demo@mangahub.dev  | demo1234  |

---

## 🌐 API Reference

### Auth
| Method | Path | Description |
|--------|------|-------------|
| POST | /api/auth/register | Create account |
| POST | /api/auth/login | Get JWT token |
| POST | /api/auth/logout | Clear session |
| GET  | /api/auth/me | Current user info |
| PUT  | /api/auth/me | Update profile |
| POST | /api/auth/change-password | Change password |

### Manga
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/manga | List/search manga |
| GET | /api/manga/popular | Popular manga |
| GET | /api/manga/{id} | Manga detail |

### Library
| Method | Path | Description |
|--------|------|-------------|
| GET  | /api/library | User's library |
| POST | /api/library | Add manga |
| PUT  | /api/library/{id}/progress | Update progress |
| DELETE | /api/library/{id} | Remove manga |
| GET  | /api/library/stats | Reading stats |

### Other
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/history | Reading history |
| GET | /api/chat/rooms | Chat rooms |
| GET | /api/chat/rooms/{id}/messages | Room messages |
| POST | /api/chat/rooms/{id}/messages | Send message |
| GET | /api/notifications | User notifications |
| GET | /api/admin/health | Server health |
| GET | /api/admin/ping | Ping |
| GET | /api/admin/stats | Server statistics |
| GET | /api/admin/logs | Server logs |

### WebSocket
Connect to `/ws?token=<JWT>` to receive:
- `connected` — welcome event
- `chat_message` — new chat message
- `notification` — new notification
- `progress_broadcast` — (via TCP relay) reading progress

Client sends:
- `{"type":"join_room","payload":{"room_id":"<id>"}}` — join a chat room
- `{"type":"ping"}` — keep-alive

### TCP (port 9001)
Line-delimited JSON messages:
```json
{"type":"identify","payload":"<user_id>"}
{"type":"progress_update","payload":{"manga_id":"...","chapter":42}}
{"type":"ping"}
```

### UDP (port 9002)
Datagram JSON:
```json
{"type":"subscribe","topics":["chapter_release","friend_activity"]}
{"type":"unsubscribe"}
{"type":"ping"}
```

### gRPC (port 9003)
Service: `manga.MangaService`
- `GetManga(id)` → MangaResponse
- `SearchManga(query, status, genre, page, limit)` → MangaListResponse
- `GetUserStats(user_id)` → UserStatsResponse
- `Ping(message)` → PongResponse

Use [evans](https://github.com/ktr0731/evans) or [grpcurl](https://github.com/fullstorydev/grpcurl) to test:
```bash
grpcurl -plaintext localhost:9003 manga.MangaService/Ping
```

---

## 🗄️ Database Schema

14 tables across 7 migration files:
1. `users` — accounts and profiles
2. `sessions` — JWT session tracking
3. `manga` — catalog with genres[], full-text indexes
4. `library_entries` — user reading lists + progress
5. `reading_history` — chapter read log
6. `reviews` — ratings and reviews
7. `friends` — follow/friend graph
8. `activity_feed` — social activity events
9. `chat_rooms` — general and manga-specific rooms
10. `chat_messages` — chat history
11. `notifications` — in-app notifications
12. `user_settings` — user preferences
13. `server_logs` — server event log
14. `backups` — backup records

---

## 🧪 Testing

```powershell
& "C:\Program Files\Go\bin\go.exe" test ./... -v
```

Test specific packages:
```powershell
& "C:\Program Files\Go\bin\go.exe" test ./internal/auth/... -v
```

---

## 🎭 Demo Checklist

- [ ] Open http://localhost:8080 → landing page
- [ ] Click "Get Started" → register new account
- [ ] Login with demo@mangahub.dev / demo1234
- [ ] Browse manga (35 pre-loaded titles)
- [ ] Click a manga → view details → "Add to Library"
- [ ] Go to Library → update reading progress
- [ ] Open Chat → join General room → send a message
- [ ] Open Admin panel → view health, ping, logs
- [ ] Open Statistics → export library JSON
- [ ] Open two browser tabs → see real-time WebSocket connection
- [ ] Check /api/admin/health for all 5 service statuses

---

## 🔧 Troubleshooting

**Port already in use:**
```powershell
netstat -ano | findstr :8080
taskkill /PID <pid> /F
```

**Database connection failed:**
- Check PostgreSQL is running: `pg_ctl status`
- Verify `.env` credentials match your PostgreSQL setup
- Ensure the database exists: `psql -U mangahub -d mangahub_db`

**Migrations failed:**
- Check PostgreSQL logs
- Ensure `pg_trgm` extension is available (included in PostgreSQL contrib)

**WebSocket not connecting:**
- Ensure you're logged in (JWT in localStorage)
- Check browser console for CORS errors

---

## 📡 Protocol Architecture

```
Browser  ────HTTP/REST──────►  :8080  (Go HTTP server)
Browser  ────WebSocket──────►  :8080/ws  (gorilla/websocket hub)
Client   ────TCP────────────►  :9001  (progress sync)
Client   ────UDP────────────►  :9002  (chapter notifications)
Internal ────gRPC───────────►  :9003  (manga/stats service)
All      ────PostgreSQL─────►  :5432  (data persistence)
```

---

## 📝 Report Summary

### What is fully implemented:
- Complete Go backend with modular architecture
- PostgreSQL with 14 tables, migrations, and seed data
- JWT authentication (register, login, logout, profile, password change)
- Manga catalog (35 titles, search, filter, detail)
- Library management (add, update, progress, remove)
- Reading history tracking
- WebSocket real-time hub (chat, notifications, presence)
- TCP sync server (progress broadcast)
- UDP notification server (chapter releases, subscriptions)
- gRPC service (MangaService with 4 RPC methods)
- Admin panel (health, stats, logs, ping, backup info)
- 15 complete web pages with dark-mode UI
- Password strength indicator, toast notifications, skeleton loaders
- Data export (library + history to JSON)

### What is scaffolded (ready to extend):
- Social/friends/activity feed (DB tables exist, UI hooks in place)
- Reviews (DB table exists, UI detail page placeholders present)
- Notification preferences (DB row exists)
- Full gRPC streaming (server supports reflection/evans)
- TCP client reconnection logic
- UDP client implementation

### AI usage disclosure:
This project was generated with AI assistance (Antigravity/Claude) for the MangaHub Net-Centric Programming course project. All architecture decisions, protocol choices, schema design, and implementation details were reviewed and guided by the student team.
