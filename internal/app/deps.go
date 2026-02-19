package app

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http"
	"github.com/ap1-final-mini-moodle/internal/domain/categorymanager"
	adminRepo "github.com/ap1-final-mini-moodle/internal/repository/admin"
	announcementRepo "github.com/ap1-final-mini-moodle/internal/repository/announcement"
	appealRepo "github.com/ap1-final-mini-moodle/internal/repository/appeal"
	assignmentRepo "github.com/ap1-final-mini-moodle/internal/repository/assignment"
	attendanceRepo "github.com/ap1-final-mini-moodle/internal/repository/attendance"
	authRepo "github.com/ap1-final-mini-moodle/internal/repository/auth"
	categorymanagerRepo "github.com/ap1-final-mini-moodle/internal/repository/categorymanager"
	chatRepo "github.com/ap1-final-mini-moodle/internal/repository/chat"
	courseRepo "github.com/ap1-final-mini-moodle/internal/repository/course"
	coursecategoryRepo "github.com/ap1-final-mini-moodle/internal/repository/coursecategory"
	dashboardRepo "github.com/ap1-final-mini-moodle/internal/repository/dashboard"
	enrollmentRepo "github.com/ap1-final-mini-moodle/internal/repository/enrollment"
	gradeRepo "github.com/ap1-final-mini-moodle/internal/repository/grade"
	groupRepo "github.com/ap1-final-mini-moodle/internal/repository/group"
	managerRepo "github.com/ap1-final-mini-moodle/internal/repository/manager"
	notificationRepo "github.com/ap1-final-mini-moodle/internal/repository/notification"
	quizRepo "github.com/ap1-final-mini-moodle/internal/repository/quiz"
	studentRepo "github.com/ap1-final-mini-moodle/internal/repository/student"
	submissionRepo "github.com/ap1-final-mini-moodle/internal/repository/submission"
	teacherRepo "github.com/ap1-final-mini-moodle/internal/repository/teacher"
	uploadRepo "github.com/ap1-final-mini-moodle/internal/repository/upload"
	userRepo "github.com/ap1-final-mini-moodle/internal/repository/user"
	"github.com/ap1-final-mini-moodle/internal/shared/config"
	"github.com/ap1-final-mini-moodle/internal/shared/database"

	adminUC "github.com/ap1-final-mini-moodle/internal/usecase/admin"
	announcementUC "github.com/ap1-final-mini-moodle/internal/usecase/announcement"
	appealUC "github.com/ap1-final-mini-moodle/internal/usecase/appeal"
	assignmentUC "github.com/ap1-final-mini-moodle/internal/usecase/assignment"
	attendanceUC "github.com/ap1-final-mini-moodle/internal/usecase/attendance"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	categorymanagerUC "github.com/ap1-final-mini-moodle/internal/usecase/categorymanager"
	chatUC "github.com/ap1-final-mini-moodle/internal/usecase/chat"
	courseUC "github.com/ap1-final-mini-moodle/internal/usecase/course"
	coursecategoryUC "github.com/ap1-final-mini-moodle/internal/usecase/coursecategory"
	dashboardUC "github.com/ap1-final-mini-moodle/internal/usecase/dashboard"
	enrollmentUC "github.com/ap1-final-mini-moodle/internal/usecase/enrollment"
	gradeUC "github.com/ap1-final-mini-moodle/internal/usecase/grade"
	groupUC "github.com/ap1-final-mini-moodle/internal/usecase/group"
	managerUC "github.com/ap1-final-mini-moodle/internal/usecase/manager"
	notificationUC "github.com/ap1-final-mini-moodle/internal/usecase/notification"
	quizUC "github.com/ap1-final-mini-moodle/internal/usecase/quiz"
	studentUC "github.com/ap1-final-mini-moodle/internal/usecase/student"
	submissionUC "github.com/ap1-final-mini-moodle/internal/usecase/submission"
	teacherUC "github.com/ap1-final-mini-moodle/internal/usecase/teacher"
	uploadUC "github.com/ap1-final-mini-moodle/internal/usecase/upload"
	userUC "github.com/ap1-final-mini-moodle/internal/usecase/user"

	authShared "github.com/ap1-final-mini-moodle/internal/shared/auth"
	"github.com/ap1-final-mini-moodle/internal/shared/websocket"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	DB *pgxpool.Pool

	UserSvc            *userUC.Service
	CourseSvc          *courseUC.Service
	EnrollmentSvc      *enrollmentUC.Service
	AssignmentSvc      *assignmentUC.Service
	SubmissionSvc      *submissionUC.Service
	GradeSvc           *gradeUC.Service
	AttendanceSvc      *attendanceUC.Service
	NotificationSvc    *notificationUC.Service
	ChatSvc            *chatUC.Service
	WSHub              *websocket.Hub
	UploadSvc          *uploadUC.Service
	GroupSvc           *groupUC.Service
	AnnouncementSvc    *announcementUC.Service
	QuizSvc            *quizUC.Service
	AppealSvc          *appealUC.Service
	DashboardSvc       *dashboardUC.Service
	StudentSvc         *studentUC.Service
	TeacherSvc         *teacherUC.Service
	AdminSvc           *adminUC.Service
	ManagerSvc         *managerUC.Service
	AuthSvc            *authUC.Service
	CategoryManagerSvc *categorymanagerUC.Service
	CourseCategorySvc  *coursecategoryUC.Service

	// Repository for middleware access
	CategoryManagerRepo categorymanager.Repository
}

func BuildDeps(db *pgxpool.Pool, redis *database.RedisClient, cfg *config.Config) *Deps {
	// Repositories
	userRepository := userRepo.NewPostgresRepository(db)
	courseRepository := courseRepo.NewRepository(db, redis)
	enrollmentRepository := enrollmentRepo.NewPostgresRepository(db)
	assignmentRepository := assignmentRepo.NewPostgresRepository(db)
	submissionRepository := submissionRepo.NewPostgresRepository(db)
	gradeRepository := gradeRepo.NewPostgresRepository(db)
	attendanceRepository := attendanceRepo.NewPostgresRepository(db)
	notificationRepository := notificationRepo.NewPostgresRepository(db)
	chatRepository := chatRepo.NewPostgresRepository(db)
	uploadRepository := uploadRepo.NewPostgresRepository(db)
	groupRepository := groupRepo.NewPostgresRepository(db)
	announcementRepository := announcementRepo.NewPostgresRepository(db)
	quizRepository := quizRepo.NewPostgresRepository(db)
	appealRepository := appealRepo.NewPostgresRepository(db)
	dashboardRepository := dashboardRepo.NewPostgresRepository(db, redis)
	studentRepository := studentRepo.NewPostgresRepository(db)
	teacherRepository := teacherRepo.NewPostgresRepository(db)
	adminRepository := adminRepo.NewPostgresRepository(db)
	managerRepository := managerRepo.NewPostgresRepository(db)
	categorymanagerRepository := categorymanagerRepo.NewPostgresRepository(db)
	coursecategoryRepository := coursecategoryRepo.NewPostgresRepository(db)
	authRepository := authRepo.NewPostgresRepository(db, redis)

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

	// WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()
	uploadService := uploadUC.NewService(uploadRepository)
	groupService := groupUC.NewService(groupRepository)
	announcementService := announcementUC.NewService(announcementRepository)
	quizService := quizUC.NewService(quizRepository)
	appealService := appealUC.NewService(appealRepository)
	dashboardService := dashboardUC.NewService(dashboardRepository)
	studentService := studentUC.NewService(studentRepository, userRepository)
	teacherService := teacherUC.NewService(teacherRepository, userRepository)
	adminService := adminUC.NewService(adminRepository, userRepository)
	managerService := managerUC.NewService(managerRepository, userRepository)
	categorymanagerService := categorymanagerUC.NewService(categorymanagerRepository, userRepository)
	coursecategoryService := coursecategoryUC.NewService(coursecategoryRepository)

	jwtConfig := &authShared.JWTConfig{
		SecretKey:            []byte(cfg.JWTSecret),
		AccessTokenDuration:  cfg.AccessTokenDuration,
		RefreshTokenDuration: cfg.RefreshTokenDuration,
		Issuer:               "mini-moodle",
		Audience:             "mini-moodle-api",
	}
	authService := authUC.NewService(authRepository, userRepository, jwtConfig)

	return &Deps{
		DB:                  db,
		UserSvc:             userService,
		CourseSvc:           courseService,
		EnrollmentSvc:       enrollmentService,
		AssignmentSvc:       assignmentService,
		SubmissionSvc:       submissionService,
		GradeSvc:            gradeService,
		AttendanceSvc:       attendanceService,
		NotificationSvc:     notificationService,
		ChatSvc:             chatService,
		WSHub:               wsHub,
		UploadSvc:           uploadService,
		GroupSvc:            groupService,
		AnnouncementSvc:     announcementService,
		QuizSvc:             quizService,
		AppealSvc:           appealService,
		DashboardSvc:        dashboardService,
		StudentSvc:          studentService,
		TeacherSvc:          teacherService,
		AdminSvc:            adminService,
		ManagerSvc:          managerService,
		AuthSvc:             authService,
		CategoryManagerSvc:  categorymanagerService,
		CourseCategorySvc:   coursecategoryService,
		CategoryManagerRepo: categorymanagerRepository,
	}
}

func BuildHTTPModules(d *Deps, jwtSecret string) []http.RoutesRegistrar {
	// Handlers
	userHandler := http.NewUserHandler(d.UserSvc)
	courseHandler := http.NewCourseHandler(d.CourseSvc)
	assignmentHandler := http.NewAssignmentHandler(d.AssignmentSvc, d.NotificationSvc, d.EnrollmentSvc)
	enrollmentHandler := http.NewEnrollmentHandler(d.EnrollmentSvc)
	submissionHandler := http.NewSubmissionHandler(d.SubmissionSvc)
	gradeHandler := http.NewGradeHandler(d.GradeSvc, d.NotificationSvc, d.SubmissionSvc)
	attendanceHandler := http.NewAttendanceHandler(d.AttendanceSvc, d.NotificationSvc)
	notificationHandler := http.NewNotificationHandler(d.NotificationSvc)
	chatHandler := http.NewChatHandler(d.ChatSvc, d.AuthSvc, d.WSHub)
	uploadHandler := http.NewUploadHandler(d.UploadSvc)
	groupHandler := http.NewGroupHandler(d.GroupSvc)
	announcementHandler := http.NewAnnouncementHandler(d.AnnouncementSvc)
	quizHandler := http.NewQuizHandler(d.QuizSvc)
	appealHandler := http.NewAppealHandler(d.AppealSvc)
	dashboardHandler := http.NewDashboardHandler(d.DashboardSvc)
	studentHandler := http.NewStudentHandler(d.StudentSvc)
	teacherHandler := http.NewTeacherHandler(d.TeacherSvc)
	adminHandler := http.NewAdminHandler(d.AdminSvc)
	managerHandler := http.NewManagerHandler(d.ManagerSvc)
	categorymanagerHandler := http.NewCategoryManagerHandler(d.CategoryManagerSvc)
	coursecategoryHandler := http.NewCourseCategoryHandler(d.CourseCategorySvc)
	authHandler := http.NewAuthHandler(d.UserSvc, d.AuthSvc)

	// Modules (all protected by JWT via authService)
	return []http.RoutesRegistrar{
		http.NewAuthModule(authHandler, d.AuthSvc),
		http.NewUserModule(userHandler, d.AuthSvc),
		http.NewGroupModule(groupHandler, d.AuthSvc),
		http.NewCourseModule(courseHandler, d.AuthSvc),
		http.NewAssignmentModule(assignmentHandler, d.AuthSvc),
		http.NewEnrollmentModule(enrollmentHandler, d.AuthSvc),
		http.NewSubmissionModule(submissionHandler, d.AuthSvc),
		http.NewGradeModule(gradeHandler, d.AuthSvc),
		http.NewAttendanceModule(attendanceHandler, d.AuthSvc),
		http.NewNotificationModule(notificationHandler, d.AuthSvc),
		http.NewChatModule(chatHandler, d.AuthSvc),
		http.NewUploadModule(uploadHandler, d.AuthSvc),
		http.NewAnnouncementModule(announcementHandler, d.AuthSvc),
		http.NewQuizModule(quizHandler, d.AuthSvc),
		http.NewAppealModule(appealHandler, d.AuthSvc),
		http.NewDashboardModule(dashboardHandler, d.AuthSvc),
		http.NewStudentModule(studentHandler, d.AuthSvc),
		http.NewTeacherModule(teacherHandler, d.AuthSvc),
		http.NewAdminModule(adminHandler, d.AuthSvc),
		http.NewManagerModule(managerHandler, d.AuthSvc),
		http.NewCategoryManagerModule(categorymanagerHandler, d.AuthSvc, d.CategoryManagerRepo),
		http.NewCourseCategoryModule(coursecategoryHandler, d.AuthSvc),
	}
}
