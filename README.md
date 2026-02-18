# Mini Moodle LMS

Learning Management System built with Go + Next.js.

## Quick Start

### Demo Credentials

```
Email: john@example.edu.kz
Password: John0102
```

### Run with Docker

```bash
make docker-up
make migrate-up
```

API: `http://localhost:8080/api/v1`  
Frontend: `http://localhost:3000`

### Local Development

```bash
# Backend
go mod download
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/moodle"
make run

# Frontend
cd frontend && npm install && npm run dev
```

## Commands

| Command | Description |
|---------|-------------|
| `make run` | Start API server |
| `make docker-up` | Start all containers |
| `make docker-down` | Stop containers |
| `make migrate-up` | Run migrations |
| `make test` | Run tests |
| `make docs` | Generate Swagger docs |

## Tech Stack

- **Backend**: Go 1.25, Gin, PostgreSQL, JWT
- **Frontend**: Next.js 14, TypeScript, Tailwind CSS, shadcn/ui
- **Infrastructure**: Docker, docker-compose

## API Endpoints

All endpoints require `Authorization: Bearer <JWT_TOKEN>`

| Resource | Endpoints |
|----------|-----------|
| Auth | `/auth/login`, `/auth/refresh`, `/auth/logout` |
| Users | `/users`, `/users/{id}` |
| Students | `/students`, `/students/{id}` |
| Courses | `/courses`, `/courses/{id}` |
| Assignments | `/assignments`, `/assignments/{id}` |
| Submissions | `/submissions`, `/submissions/{id}` |
| Grades | `/grades`, `/grades/{id}` |
| Attendance | `/attendance` |
| Chat | `/messages`, WebSocket `/chat` |
| Groups | `/groups`, `/groups/{id}` |

Full API documentation: `docs/swagger.yaml`

## Environment Variables

| Variable | Required | Default |
|----------|----------|---------|
| `DATABASE_URL` | Yes | - |
| `PORT` | No | 8080 |
| `JWT_SECRET` | No | - |

## License

Private - Educational use only.
