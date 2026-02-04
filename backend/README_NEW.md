# Mini-Moodle LMS - Компактный проект с Clean Architecture

Полнофункциональная LMS система с чистой архитектурой, разработанная на Go.

## 📋 Быстрый старт

### Требования
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+ (опционально)

### Установка и запуск

```bash
# 1. Клонировать проект
git clone <repo>
cd ap1-final-mini-moodle/backend

# 2. Запустить с Docker (самый простой способ)
make docker-up

# 3. Применить миграции БД
make migrate-up

# 4. API доступен на http://localhost:8080/api/v1

# Для остановки
make docker-down
```

### Локальный запуск (без Docker)

```bash
# Требуется запущенный PostgreSQL
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/moodle"

# Запустить сервер
make run

# Или напрямую
go run cmd/api/main.go
```

## 🏗️ Архитектура

### Clean Architecture (4 слоя)

```
┌─────────────────────────────────────┐
│   HTTP Handlers (Gin)               │ ← REST API
├─────────────────────────────────────┤
│   Usecase/Service (Business Logic)  │ ← Бизнес-логика
├─────────────────────────────────────┤
│   Repository (Data Access)          │ ← PostgreSQL/Redis
├─────────────────────────────────────┤
│   Domain (Entities + Interfaces)    │ ← Ядро системы
└─────────────────────────────────────┘
```

### Слои проекта

| Слой | Путь | Назначение |
|------|------|-----------|
| **Domain** | `internal/domain/*/` | Entities, interfaces, бизнес-правила |
| **Repository** | `internal/repository/*/` | PostgreSQL implementation |
| **Usecase** | `internal/usecase/*/` | Service layer, бизнес-логика |
| **Delivery** | `internal/delivery/http/` | HTTP handlers (Gin) |

## 📦 Домены (10 сервисов)

### Core
- **User** - Аутентификация и управление пользователями (4 роли: admin, teacher, student, manager)
- **Course** - Создание и управление курсами
- **Enrollment** - Регистрация студентов на курсы

### Academic
- **Assignment** - Создание заданий
- **Submission** - Сдача студентами работ
- **Grade** - Выставление оценок
- **Attendance** - Учет посещаемости

### Communication
- **Chat** - Обсуждение в курсе
- **Notification** - Уведомления
- **Upload** - Загрузка файлов

## 🔌 API Endpoints

Все endpoints под `/api/v1/`

### Users
```
POST   /users                    Регистрация
GET    /users                    Список пользователей
GET    /users/{id}               Получить пользователя
PATCH  /users/{id}               Обновить пользователя
DELETE /users/{id}               Удалить пользователя
```

### Courses
```
POST   /courses                  Создать курс
GET    /courses                  Список курсов
GET    /courses/{id}             Получить курс
PATCH  /courses/{id}             Обновить курс
DELETE /courses/{id}             Удалить курс
```

### Enrollments
```
POST   /enrollments              Записать студента
GET    /enrollments/course/{id}  По курсу
GET    /enrollments/student/{id} По студенту
DELETE /enrollments/{id}         Отписать
```

### Similar pattern для assignments, submissions, grades, attendance, chat, notifications, uploads

## 🛠️ Инструменты и команды

### Основные команды

```bash
make build         # Собрать бинарник
make run           # Запустить локально
make test          # Запустить тесты
make clean         # Очистить артефакты
```

### Docker команды

```bash
make docker-up     # Поднять контейнеры
make docker-down   # Опустить контейнеры
make docker-logs   # Смотреть логи API
make docker-build  # Пересобрать образ
```

### Database команды

```bash
make migrate-up    # Применить миграции
make migrate-down  # Откатить миграции
```

## 📝 Примеры использования

### Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "role": "student"
  }'
```

### Создание курса

```bash
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Content-Type: application/json" \
  -d '{
    "code": "CS101",
    "title": "Introduction to Computer Science",
    "description": "Basic CS course",
    "teacher_id": "teacher_id_here",
    "max_points": 100
  }'
```

### Запись студента на курс

```bash
curl -X POST http://localhost:8080/api/v1/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": "course_id_here",
    "student_id": "student_id_here"
  }'
```

## 🧪 Тестирование

```bash
# Запустить все тесты
make test

# Запустить с coverage
go test -v ./... -cover

# Unit тесты
go test ./tests/unit/...

# Integration тесты
go test ./tests/integration/...
```

## 📚 Документация

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Подробно о архитектуре
- [REFACTORING_SUMMARY.md](./REFACTORING_SUMMARY.md) - Что было рефакторено
- [docs/swagger.yaml](./docs/swagger.yaml) - API specification

## 🐳 Docker Compose

### Сервисы

```yaml
- postgres:5432  (база данных)
- redis:6379     (кеш и сессии)
- api:8080       (приложение)
```

### Environment Variables

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=moodle
REDIS_URL=redis://redis:6379
PORT=8080
```

### Volumes

```bash
postgres_data/  - Данные PostgreSQL
redis_data/     - Данные Redis
```

## 📊 Производительность

- **Multi-stage Docker** - 120MB финальный образ
- **Connection pooling** - pgxpool для эффективности
- **Paginated queries** - Избегаем huge response
- **Indexed columns** - Быстрые запросы

## 🔐 Безопасность

Текущая версия для обучения. Для продакшена добавьте:
- JWT authentication
- CORS configuration
- Input validation
- Rate limiting
- HTTPS/TLS

## 🚀 Деплой

### На Ubuntu/VPS

```bash
# 1. Установить Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 2. Клонировать и запустить
git clone <repo> && cd ap1-final-mini-moodle/backend
docker-compose up -d

# 3. Проверить статус
docker-compose ps
```

### На Kubernetes

```bash
# Использовать Dockerfile для создания образа
docker build -t lms-api:latest .
docker tag lms-api:latest your-registry/lms-api:latest
docker push your-registry/lms-api:latest

# Развернуть
kubectl apply -f deployment.yaml
```

## 📈 Структура проекта

```
backend/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── domain/                  # 10 domain packages
│   ├── repository/              # 10 repository implementations
│   ├── usecase/                 # 10 service implementations
│   ├── delivery/http/           # HTTP handlers
│   └── shared/                  # Utilities
├── migrations/                  # Database migrations
├── tests/                       # Unit & integration tests
├── Dockerfile                   # Container image
├── docker-compose.yml          # Local dev setup
├── Makefile                    # Build commands
└── go.mod                      # Dependencies
```

## 💻 Системные требования

| Компонент | Требование |
|-----------|-----------|
| Go | 1.21+ |
| PostgreSQL | 15+ |
| Redis | 7+ (опционально) |
| Docker | 20+ (опционально) |
| RAM | 512MB минимум |
| Disk | 5GB для БД |

## 📞 Troubleshooting

### Проблема: Connection refused на БД

```bash
# Решение: Убедиться что PostgreSQL запущен
docker-compose ps

# Или проверить локальный PostgreSQL
psql -U postgres -c "SELECT version();"
```

### Проблема: Port 8080 уже занят

```bash
# Решение: Изменить PORT
PORT=3000 make run

# Или в docker-compose.yml
ports:
  - "3000:8080"
```

### Проблема: Database не инициализируется

```bash
# Решение: Применить миграции
make migrate-up

# Или вручную
go run cmd/migrate/main.go up
```

## 📖 Дополнительные ресурсы

- [Clean Architecture by Robert Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Gin Documentation](https://gin-gonic.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)

## 📝 Лицензия

Приватная - Только для образовательных целей

---

**Версия**: 1.0.0  
**Последнее обновление**: 2026-02-04  
**Автор**: Development Team  
**Контакт**: support@example.com
