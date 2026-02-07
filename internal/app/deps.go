package app

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http"
	announcementRepo "github.com/ap1-final-mini-moodle/internal/repository/announcement"
	appealRepo "github.com/ap1-final-mini-moodle/internal/repository/appeal"
	assignmentRepo "github.com/ap1-final-mini-moodle/internal/repository/assignment"
	attendanceRepo "github.com/ap1-final-mini-moodle/internal/repository/attendance"
	chatRepo "github.com/ap1-final-mini-moodle/internal/repository/chat"
	courseRepo "github.com/ap1-final-mini-moodle/internal/repository/course"
	dashboardRepo "github.com/ap1-final-mini-moodle/internal/repository/dashboard"
	enrollmentRepo "github.com/ap1-final-mini-moodle/internal/repository/enrollment"
	gradeRepo "github.com/ap1-final-mini-moodle/internal/repository/grade"
	groupRepo "github.com/ap1-final-mini-moodle/internal/repository/group"
	notificationRepo "github.com/ap1-final-mini-moodle/internal/repository/notification"
	quizRepo "github.com/ap1-final-mini-moodle/internal/repository/quiz"
	studentRepo "github.com/ap1-final-mini-moodle/internal/repository/student"
	submissionRepo "github.com/ap1-final-mini-moodle/internal/repository/submission"
	teacherRepo "github.com/ap1-final-mini-moodle/internal/repository/teacher"
	uploadRepo "github.com/ap1-final-mini-moodle/internal/repository/upload"
	userRepo "github.com/ap1-final-mini-moodle/internal/repository/user"

	announcementUC "github.com/ap1-final-mini-moodle/internal/usecase/announcement"
	appealUC "github.com/ap1-final-mini-moodle/internal/usecase/appeal"
	assignmentUC "github.com/ap1-final-mini-moodle/internal/usecase/assignment"
	attendanceUC "github.com/ap1-final-mini-moodle/internal/usecase/attendance"
	chatUC "github.com/ap1-final-mini-moodle/internal/usecase/chat"
	courseUC "github.com/ap1-final-mini-moodle/internal/usecase/course"
	dashboardUC "github.com/ap1-final-mini-moodle/internal/usecase/dashboard"
	enrollmentUC "github.com/ap1-final-mini-moodle/internal/usecase/enrollment"
	gradeUC "github.com/ap1-final-mini-moodle/internal/usecase/grade"
	groupUC "github.com/ap1-final-mini-moodle/internal/usecase/group"
	notificationUC "github.com/ap1-final-mini-moodle/internal/usecase/notification"
	quizUC "github.com/ap1-final-mini-moodle/internal/usecase/quiz"
	studentUC "github.com/ap1-final-mini-moodle/internal/usecase/student"
	submissionUC "github.com/ap1-final-mini-moodle/internal/usecase/submission"
	teacherUC "github.com/ap1-final-mini-moodle/internal/usecase/teacher"
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
	GroupSvc        *groupUC.Service
	AnnouncementSvc *announcementUC.Service
	QuizSvc         *quizUC.Service
	AppealSvc       *appealUC.Service
	DashboardSvc    *dashboardUC.Service
	StudentSvc      *studentUC.Service
	TeacherSvc      *teacherUC.Service
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
	groupRepository := groupRepo.NewPostgresRepository(db)
	announcementRepository := announcementRepo.NewPostgresRepository(db)
	quizRepository := quizRepo.NewPostgresRepository(db)
	appealRepository := appealRepo.NewPostgresRepository(db)
	dashboardRepository := dashboardRepo.NewPostgresRepository(db)
	studentRepository := studentRepo.NewPostgresRepository(db)
	teacherRepository := teacherRepo.NewPostgresRepository(db)

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
	groupService := groupUC.NewService(groupRepository)
	announcementService := announcementUC.NewService(announcementRepository)
	quizService := quizUC.NewService(quizRepository)
	appealService := appealUC.NewService(appealRepository)
	dashboardService := dashboardUC.NewService(dashboardRepository)
	studentService := studentUC.NewService(studentRepository)
	teacherService := teacherUC.NewService(teacherRepository)

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
		GroupSvc:        groupService,
		AnnouncementSvc: announcementService,
		QuizSvc:         quizService,
		AppealSvc:       appealService,
		DashboardSvc:    dashboardService,
		StudentSvc:      studentService,
		TeacherSvc:      teacherService,
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
	groupHandler := http.NewGroupHandler(d.GroupSvc)
	announcementHandler := http.NewAnnouncementHandler(d.AnnouncementSvc)
	quizHandler := http.NewQuizHandler(d.QuizSvc)
	appealHandler := http.NewAppealHandler(d.AppealSvc)
	dashboardHandler := http.NewDashboardHandler(d.DashboardSvc)
	studentHandler := http.NewStudentHandler(d.StudentSvc)
	teacherHandler := http.NewTeacherHandler(d.TeacherSvc)

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
		http.NewGroupModule(groupHandler, secret),
		http.NewAnnouncementModule(announcementHandler, secret),
		http.NewQuizModule(quizHandler, secret),
		http.NewAppealModule(appealHandler, secret),
		http.NewDashboardModule(dashboardHandler, secret),
		http.NewStudentModule(studentHandler, secret),
		http.NewTeacherModule(teacherHandler, secret),
	}
}
