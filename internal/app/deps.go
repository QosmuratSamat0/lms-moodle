package app

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http"
	assignmentRepo "github.com/ap1-final-mini-moodle/internal/repository/assignment"
	attendanceRepo "github.com/ap1-final-mini-moodle/internal/repository/attendance"
	chatRepo "github.com/ap1-final-mini-moodle/internal/repository/chat"
	courseRepo "github.com/ap1-final-mini-moodle/internal/repository/course"
	enrollmentRepo "github.com/ap1-final-mini-moodle/internal/repository/enrollment"
	gradeRepo "github.com/ap1-final-mini-moodle/internal/repository/grade"
	notificationRepo "github.com/ap1-final-mini-moodle/internal/repository/notification"
	submissionRepo "github.com/ap1-final-mini-moodle/internal/repository/submission"
	uploadRepo "github.com/ap1-final-mini-moodle/internal/repository/upload"
	userRepo "github.com/ap1-final-mini-moodle/internal/repository/user"

	assignmentUC "github.com/ap1-final-mini-moodle/internal/usecase/assignment"
	attendanceUC "github.com/ap1-final-mini-moodle/internal/usecase/attendance"
	chatUC "github.com/ap1-final-mini-moodle/internal/usecase/chat"
	courseUC "github.com/ap1-final-mini-moodle/internal/usecase/course"
	enrollmentUC "github.com/ap1-final-mini-moodle/internal/usecase/enrollment"
	gradeUC "github.com/ap1-final-mini-moodle/internal/usecase/grade"
	notificationUC "github.com/ap1-final-mini-moodle/internal/usecase/notification"
	submissionUC "github.com/ap1-final-mini-moodle/internal/usecase/submission"
	uploadUC "github.com/ap1-final-mini-moodle/internal/usecase/upload"
	userUC "github.com/ap1-final-mini-moodle/internal/usecase/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	DB *pgxpool.Pool

	UserSvc         *userUC.Service
	CourseSvc       *courseUC.Service
	EnrollmentSvc   *enrollmentUC.Service
	AssignmentSvc   *assignmentUC.Service
	SubmissionSvc   *submissionUC.Service
	GradeSvc        *gradeUC.Service
	AttendanceSvc   *attendanceUC.Service
	NotificationSvc *notificationUC.Service
	ChatSvc         *chatUC.Service
	UploadSvc       *uploadUC.Service
}

func BuildDeps(db *pgxpool.Pool) *Deps {
	// Repositories
	userRepository := userRepo.NewPostgresRepository(db)
	courseRepository := courseRepo.NewPostgresRepository(db)
	enrollmentRepository := enrollmentRepo.NewPostgresRepository(db)
	assignmentRepository := assignmentRepo.NewPostgresRepository(db)
	submissionRepository := submissionRepo.NewPostgresRepository(db)
	gradeRepository := gradeRepo.NewPostgresRepository(db)
	attendanceRepository := attendanceRepo.NewPostgresRepository(db)
	notificationRepository := notificationRepo.NewPostgresRepository(db)
	chatRepository := chatRepo.NewPostgresRepository(db)
	uploadRepository := uploadRepo.NewPostgresRepository(db)

	// Services
	userService := userUC.NewService(userRepository)
	courseService := courseUC.NewService(courseRepository)
	enrollmentService := enrollmentUC.NewService(enrollmentRepository)
	assignmentService := assignmentUC.NewService(assignmentRepository)
	submissionService := submissionUC.NewService(submissionRepository)
	gradeService := gradeUC.NewService(gradeRepository)
	attendanceService := attendanceUC.NewService(attendanceRepository)
	notificationService := notificationUC.NewService(notificationRepository)
	chatService := chatUC.NewService(chatRepository)
	uploadService := uploadUC.NewService(uploadRepository)

	return &Deps{
		DB:              db,
		UserSvc:         userService,
		CourseSvc:       courseService,
		EnrollmentSvc:   enrollmentService,
		AssignmentSvc:   assignmentService,
		SubmissionSvc:   submissionService,
		GradeSvc:        gradeService,
		AttendanceSvc:   attendanceService,
		NotificationSvc: notificationService,
		ChatSvc:         chatService,
		UploadSvc:       uploadService,
	}
}

func BuildHTTPModules(d *Deps, jwtSecret string) []http.RoutesRegistrar {
	secret := []byte(jwtSecret)

	// Handlers
	userHandler := http.NewUserHandler(d.UserSvc)
	courseHandler := http.NewCourseHandler(d.CourseSvc)
	assignmentHandler := http.NewAssignmentHandler(d.AssignmentSvc)
	enrollmentHandler := http.NewEnrollmentHandler(d.EnrollmentSvc)
	submissionHandler := http.NewSubmissionHandler(d.SubmissionSvc)
	gradeHandler := http.NewGradeHandler(d.GradeSvc)
	attendanceHandler := http.NewAttendanceHandler(d.AttendanceSvc)
	notificationHandler := http.NewNotificationHandler(d.NotificationSvc)
	chatHandler := http.NewChatHandler(d.ChatSvc)
	uploadHandler := http.NewUploadHandler(d.UploadSvc)

	// Modules
	return []http.RoutesRegistrar{
		http.NewUserModule(userHandler, secret),
		http.NewCourseModule(courseHandler, secret),
		http.NewAssignmentModule(assignmentHandler, secret),
		http.NewEnrollmentModule(enrollmentHandler, secret),
		http.NewSubmissionModule(submissionHandler, secret),
		http.NewGradeModule(gradeHandler, secret),
		http.NewAttendanceModule(attendanceHandler, secret),
		http.NewNotificationModule(notificationHandler, secret),
		http.NewChatModule(chatHandler, secret),
		http.NewUploadModule(uploadHandler, secret),
	}
}
