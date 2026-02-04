// internal/api/router.go
package api

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/api/routes"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/analytics"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/assignment"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/attendance"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/chat"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/course"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/enrollment"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/grade"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/group"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/manager"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/notification"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/plagiarism"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/schedule"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/session"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/student"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/submission"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/teacher"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/upload"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/user"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// Handlers contains all domain handlers
type Handlers struct {
	User         *user.Handler
	Course       *course.Handler
	Enrollment   *enrollment.Handler
	Assignment   *assignment.Handler
	Submission   *submission.Handler
	Grade        *grade.Handler
	Attendance   *attendance.Handler
	Chat         *chat.Handler
	Notification *notification.Handler
	Schedule     *schedule.Handler
	Analytics    *analytics.Handler
	Session      *session.Handler
	Student      *student.Handler
	Teacher      *teacher.Handler
	Manager      *manager.Handler
	Plagiarism   *plagiarism.Handler
	Upload       *upload.Handler
	Group        *group.Handler
}

// Router manages all API routes
type Router struct {
	engine   *gin.Engine
	auth     *middleware.AuthMiddleware
	handlers *Handlers
}

// NewRouter creates a new router instance
func NewRouter(engine *gin.Engine, auth *middleware.AuthMiddleware, handlers *Handlers) *Router {
	return &Router{
		engine:   engine,
		auth:     auth,
		handlers: handlers,
	}
}

// SetupRoutes configures all API routes
func (r *Router) SetupRoutes() {
	// API version prefix
	api := r.engine.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes (no auth)
	publicRouter := routes.NewPublicRouter(r.handlers.User)
	publicRouter.SetupRoutes(api)

	// Protected routes (require auth)
	protected := api.Group("")
	protected.Use(r.auth.Authenticate())

	// Initialize routers using clean architecture pattern
	userRouter := routes.NewUserRouter(r.handlers.User)
	courseRouter := routes.NewCourseRouter(r.handlers.Course)
	enrollmentRouter := routes.NewEnrollmentRouter(r.handlers.Enrollment)
	assignmentRouter := routes.NewAssignmentRouter(r.handlers.Assignment)
	submissionRouter := routes.NewSubmissionRouter(r.handlers.Submission)
	gradeRouter := routes.NewGradeRouter(r.handlers.Grade)
	attendanceRouter := routes.NewAttendanceRouter(r.handlers.Attendance)
	chatRouter := routes.NewChatRouter(r.handlers.Chat)
	notificationRouter := routes.NewNotificationRouter(r.handlers.Notification)
	scheduleRouter := routes.NewScheduleRouter(r.handlers.Schedule)
	analyticsRouter := routes.NewAnalyticsRouter(r.handlers.Analytics)
	sessionRouter := routes.NewSessionRouter(r.handlers.Session)
	studentRouter := routes.NewStudentRouter(r.handlers.Student)
	teacherRouter := routes.NewTeacherRouter(r.handlers.Teacher, r.handlers.Group)
	managerRouter := routes.NewManagerRouter(r.handlers.Manager)
	plagiarismRouter := routes.NewPlagiarismRouter(r.handlers.Plagiarism)
	uploadRouter := routes.NewUploadRouter(r.handlers.Upload)
	groupRouter := routes.NewGroupRouter(r.handlers.Group)

	// Setup all routes
	userRouter.SetupRoutes(protected)
	courseRouter.SetupRoutes(protected)
	enrollmentRouter.SetupRoutes(protected)
	assignmentRouter.SetupRoutes(protected)
	submissionRouter.SetupRoutes(protected)
	gradeRouter.SetupRoutes(protected)
	attendanceRouter.SetupRoutes(protected)
	chatRouter.SetupRoutes(protected)
	notificationRouter.SetupRoutes(protected)
	scheduleRouter.SetupRoutes(protected)
	analyticsRouter.SetupRoutes(protected)
	sessionRouter.SetupRoutes(protected)
	studentRouter.SetupRoutes(protected)
	teacherRouter.SetupRoutes(protected)
	managerRouter.SetupRoutes(protected)
	plagiarismRouter.SetupRoutes(protected)
	uploadRouter.SetupRoutes(protected)
	groupRouter.SetupRoutes(protected)
}
