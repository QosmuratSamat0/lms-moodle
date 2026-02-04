# Routes Package - Clean Architecture

## Обзор

Пакет `routes` организован согласно принципам **Clean Architecture**. Каждый роутер отвечает за свой домен и инкапсулирует логику маршрутизации.

## Структура

```
routes/
├── public.go          # Публичные маршруты (регистрация, логин)
├── user.go           # Пользовательские маршруты
├── course.go         # Маршруты курсов
├── enrollment.go     # Маршруты записи на курсы
├── assignment.go     # Маршруты заданий
├── submission.go     # Маршруты сдачи работ
├── grade.go          # Маршруты оценок
├── attendance.go     # Маршруты посещаемости
├── chat.go           # Маршруты чата
├── notification.go   # Маршруты уведомлений
├── schedule.go       # Маршруты расписания
├── analytics.go      # Маршруты аналитики
├── session.go        # Маршруты сессий
├── student.go        # Маршруты студентов
├── teacher.go        # Маршруты преподавателей
├── manager.go        # Маршруты менеджеров
├── plagiarism.go     # Маршруты проверки плагиата
├── upload.go         # Маршруты загрузки файлов
├── group.go          # Маршруты групп
└── admin.go          # Зарезервировано для будущих админ маршрутов
```

## Принципы

### 1. Инкапсуляция
Каждый роутер инкапсулирует:
- Конструктор (`New*Router`)
- Метод настройки маршрутов (`SetupRoutes`)
- Необходимые зависимости (handlers)

### 2. Разделение ответственности
- **Router** - только маршрутизация и middleware
- **Handler** - обработка HTTP запросов
- **Service** - бизнес-логика
- **Repository** - работа с данными

### 3. Независимость
Каждый роутер может быть:
- Протестирован независимо
- Изменен без влияния на другие роутеры
- Повторно использован

## Пример использования

### Создание нового роутера

```go
// routes/example.go
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/example"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

type ExampleRouter struct {
	exampleHandler *example.Handler
}

func NewExampleRouter(exampleHandler *example.Handler) *ExampleRouter {
	return &ExampleRouter{
		exampleHandler: exampleHandler,
	}
}

func (er *ExampleRouter) SetupRoutes(api *gin.RouterGroup) {
	examples := api.Group("/examples")
	{
		// Публичные маршруты
		examples.GET("", er.exampleHandler.List)
		examples.GET("/:id", er.exampleHandler.GetByID)

		// Защищенные маршруты (только для админов)
		adminRoutes := examples.Group("")
		adminRoutes.Use(middleware.RequireRole("admin"))
		{
			adminRoutes.POST("", er.exampleHandler.Create)
			adminRoutes.PUT("/:id", er.exampleHandler.Update)
			adminRoutes.DELETE("/:id", er.exampleHandler.Delete)
		}
	}
}
```

### Интеграция в router.go

```go
// В методе SetupRoutes главного Router
exampleRouter := routes.NewExampleRouter(r.handlers.Example)
exampleRouter.SetupRoutes(protected)
```

## Middleware

### Аутентификация
```go
protected := api.Group("")
protected.Use(r.auth.Authenticate())
```

### Авторизация по роли
```go
adminRoutes := group.Group("")
adminRoutes.Use(middleware.RequireRole("admin"))
```

Доступные роли:
- `admin` - администратор системы
- `manager` - менеджер
- `teacher` - преподаватель
- `student` - студент

## Специальные маршруты

### WebSocket (chat.go)
```go
// Отдельная функция для настройки WebSocket
func SetupChatWebSocket(engine *gin.Engine, auth *middleware.AuthMiddleware, 
                        service chat.Service, hub *websocket.Hub) {
	// WebSocket требует прямой доступ к engine
}
```

## Тестирование

Пример unit-теста для роутера:

```go
func TestUserRouter_SetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	
	mockHandler := &user.Handler{}
	router := routes.NewUserRouter(mockHandler)
	
	api := engine.Group("/api/v1")
	router.SetupRoutes(api)
	
	// Проверяем зарегистрированные маршруты
	routes := engine.Routes()
	assert.Contains(t, routes, "GET /api/v1/users/me")
}
```

## Миграция

При добавлении новых маршрутов:

1. Создайте новый файл в `routes/` или обновите существующий
2. Следуйте паттерну: Router struct → Constructor → SetupRoutes
3. Добавьте handler в `api.Handlers` если нужен новый
4. Инициализируйте роутер в `router.go`
5. Вызовите `SetupRoutes()` в нужном месте

## Преимущества текущей архитектуры

✅ **Модульность** - каждый роутер независим  
✅ **Тестируемость** - легко писать unit-тесты  
✅ **Масштабируемость** - просто добавлять новые роутеры  
✅ **Читаемость** - понятная структура файлов  
✅ **Поддерживаемость** - изменения локализованы  
✅ **Clean Architecture** - разделение ответственности  

## Рекомендации

1. **Один роутер = один домен** - не смешивайте разные доменные области
2. **Минимум логики** - роутеры только регистрируют маршруты
3. **Используйте middleware** - для общих задач (auth, logging, etc.)
4. **Документируйте** - комментируйте сложные маршруты
5. **Группируйте** - используйте router groups для организации
