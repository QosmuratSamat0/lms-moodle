golang-service/
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── migrate/
│   │   └── main.go
│   ├── worker/
│       └── main.go
├── docs/
│   ├── architecture.md
│   └── swagger.yaml
├── internal/
│   ├── api/
│   │   ├── routes/
│   │   │   ├── admin.go
│   │   │   ├── chat.go
│   │   │   ├── course.go
│   │   │   ├── manager.go
│   │   │   ├── public.go
│   │   │   ├── student.go
│   │   │   ├── teacher.go
│   │   │   └── user.go
│   │   ├── middleware.go
│   │   └── router.go
│   ├── domain/
│   │   ├── analytics/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── assignment/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── attendance/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── chat/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── websocket.go
│   │   ├── course/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── enrollment/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── grade/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── manager/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── notification/
│   │   │   ├── email.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── push.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── sms.go
│   │   ├── plagiarism/
│   │   │   ├── detector.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── schedule/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── session/
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── student/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── submission/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── teacher/
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── user/
│   │       ├── dto.go
│   │       ├── handler.go
│   │       ├── model.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── validator.go
│   ├── shared/
│       ├── config/
│       │   ├── config.go
│       │   └── loader.go
│       ├── constants/
│       │   ├── permissions.go
│       │   ├── roles.go
│       │   └── status.go
│       ├── database/
│       │   ├── postgres.go
│       │   ├── redis.go
│       │   └── transaction.go
│       ├── errorx/
│       │   ├── errors.go
│       │   └── handler.go
│       ├── middleware/
│       │   ├── auth.go
│       │   ├── cors.go
│       │   ├── logger.go
│       │   ├── rate_limit.go
│       │   ├── recovery.go
│       │   └── role.go
│       ├── utils/
│       │   ├── file.go
│       │   ├── hash.go
│       │   ├── jwt.go
│       │   ├── pagination.go
│       │   ├── response.go
│       │   └── validator.go
│       ├── websocket/
│           ├── client.go
│           ├── hub.go
│           └── message.go
├── migrations/
│   ├── 000001_init.down.sql
│   └── 000001_init.up.sql
├── pkg/
│   ├── cache/
│   │   └── redis.go
│   ├── email/
│   │   └── sender.go
│   ├── logger/
│   │   └── logger.go
│   ├── queue/
│   │   └── worker.go
│   ├── storage/
│       └── s3.go
├── scripts/
│   ├── backup.sh
│   └── seed.go
├── tests/
│   ├── e2e/
│   │   └── api_test.go
│   ├── integration/
│       ├── course_test.go
│       ├── student_test.go
│       └── teacher_test.go
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── go.sum
└── main.go
