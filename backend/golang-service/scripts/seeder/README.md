# LMS Data Seeder

Генератор тестовых данных для LMS через API эндпоинты.

## Возможности

- ✅ Регистрация и логин пользователей (admin, teachers, students)
- ✅ Создание курсов
- ✅ Запись студентов на курсы (enrollments)
- ✅ Создание заданий (assignments)
- ✅ Отправка работ (submissions)
- ✅ Выставление оценок (grades)
- ✅ Учёт посещаемости (attendance)
- ✅ Расписание (schedule)
- ✅ Чаты и сообщения

## Использование

### Запуск через Make

```bash
# Запуск с дефолтным URL (http://localhost:8080/api/v1)
make seed

# Запуск с кастомным URL
make seed-url URL=http://localhost:8080/api/v1
```

### Запуск напрямую

```bash
# Дефолтный URL
go run ./scripts/seeder/main.go

# С переменной окружения
API_URL=http://localhost:8080/api/v1 go run ./scripts/seeder/main.go
```

## Требования

- Запущенный API сервер
- Запущенная база данных с миграциями

## Тестовые учётные данные

После запуска seeder'а доступны следующие аккаунты:

### Admin/Manager
- Email: `admin@lms.local`
- Password: `Admin123!@#`

### Teachers
- `john.smith@lms.local` / `Teacher123!` (Computer Science)
- `maria.garcia@lms.local` / `Teacher123!` (Mathematics)
- `david.johnson@lms.local` / `Teacher123!` (Physics)
- `sarah.williams@lms.local` / `Teacher123!` (Computer Science)
- `michael.brown@lms.local` / `Teacher123!` (Engineering)

### Students
- `alice.student@lms.local` / `Student123!` (CS-101)
- `bob.student@lms.local` / `Student123!` (CS-101)
- `charlie.student@lms.local` / `Student123!` (CS-101)
- `diana.student@lms.local` / `Student123!` (CS-102)
- `evan.student@lms.local` / `Student123!` (CS-102)
- `fiona.student@lms.local` / `Student123!` (MATH-101)
- `george.student@lms.local` / `Student123!` (MATH-101)
- `hannah.student@lms.local` / `Student123!` (PHYS-101)
- `ivan.student@lms.local` / `Student123!` (PHYS-101)
- `julia.student@lms.local` / `Student123!` (ENG-101)

## Генерируемые данные

| Тип данных | Количество |
|------------|------------|
| Teachers | 5 |
| Students | 10 |
| Courses | 8 |
| Assignments | 3-5 на курс |
| Enrollments | 2-4 на студента |
| Submissions | 60-80% от возможных |
| Grades | 70-90% от submissions |
| Attendance sessions | 3-5 на курс |
| Schedule events | 4 на курс |
| Chat rooms | 1 на курс + direct chats |

## Структура кода

```
scripts/seeder/
└── main.go          # Основной файл seeder'а
    ├── Data Models  # Структуры для API responses
    ├── Sample Data  # Тестовые данные
    ├── SeederState  # Состояние во время работы
    ├── HTTP Client  # Функции для запросов
    └── Seeders      # Функции генерации данных
```

## Расширение

Чтобы добавить новые данные:

1. Добавьте шаблоны данных в секцию `Sample Data`
2. Создайте новую функцию `seedXxx()` в секции `Seeders`
3. Добавьте вызов в массив `seeders` в `main()`
