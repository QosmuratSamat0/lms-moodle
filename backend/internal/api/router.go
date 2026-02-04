// internal/api/router.go
package api

import (
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
	r.setupPublicRoutes(api)

	// Protected routes (require auth)
	protected := api.Group("")
	protected.Use(r.auth.Authenticate())

	r.setupUserRoutes(protected)
	r.setupCourseRoutes(protected)
	r.setupEnrollmentRoutes(protected)
	r.setupAssignmentRoutes(protected)
	r.setupSubmissionRoutes(protected)
	r.setupGradeRoutes(protected)
	r.setupAttendanceRoutes(protected)
	r.setupChatRoutes(protected)
	r.setupNotificationRoutes(protected)
	r.setupScheduleRoutes(protected)
	r.setupAnalyticsRoutes(protected)
	r.setupSessionRoutes(protected)
	r.setupStudentRoutes(protected)
	r.setupTeacherRoutes(protected)
	r.setupManagerRoutes(protected)
	r.setupPlagiarismRoutes(protected)
	r.setupUploadRoutes(protected)
	r.setupGroupRoutes(protected)
}

func (r *Router) setupPublicRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", r.handlers.User.Register)
		auth.POST("/login", r.handlers.User.Login)
		auth.POST("/refresh", r.handlers.User.RefreshToken)
	}
}

func (r *Router) setupUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/me", r.handlers.User.GetProfile)
		users.PUT("/me", r.handlers.User.UpdateProfile)
		users.PUT("/me/password", r.handlers.User.ChangePassword)

		// Admin only
		admin := users.Group("")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("", r.handlers.User.ListUsers)
			admin.GET("/:id", r.handlers.User.GetUser)
			admin.POST("/:id/deactivate", r.handlers.User.DeactivateUser)
			admin.POST("/:id/activate", r.handlers.User.ActivateUser)
		}
	}
}

func (r *Router) setupCourseRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses")
	{
		courses.GET("", r.handlers.Course.List)
		courses.GET("/:id", r.handlers.Course.GetByID)

		// Teacher/Admin only
		teacherRoutes := courses.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("", r.handlers.Course.Create)
			teacherRoutes.PUT("/:id", r.handlers.Course.Update)
			teacherRoutes.DELETE("/:id", r.handlers.Course.Delete)
		}
	}
}

func (r *Router) setupEnrollmentRoutes(api *gin.RouterGroup) {
	enrollments := api.Group("/enrollments")
	{
		// Student enrollment
		enrollments.POST("", r.handlers.Enrollment.Enroll)
		enrollments.DELETE("/:id", r.handlers.Enrollment.Drop)
		enrollments.GET("/my", r.handlers.Enrollment.ListMyEnrollments)

		// Course enrollments (teachers)
		api.GET("/courses/:id/enrollments", r.handlers.Enrollment.ListByCourse)

		// Approval (teacher/admin)
		teacherRoutes := enrollments.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin", "manager"))
		{
			teacherRoutes.PUT("/:id/status", r.handlers.Enrollment.UpdateStatus)
		}
	}
}

func (r *Router) setupAssignmentRoutes(api *gin.RouterGroup) {
	assignments := api.Group("/assignments")
	{
		assignments.GET("/:id", r.handlers.Assignment.GetByID)

		// Course assignments
		api.GET("/courses/:id/assignments", r.handlers.Assignment.ListByCourse)

		// Teacher/Admin only
		teacherRoutes := assignments.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("", r.handlers.Assignment.Create)
			teacherRoutes.PUT("/:id", r.handlers.Assignment.Update)
			teacherRoutes.DELETE("/:id", r.handlers.Assignment.Delete)
		}
	}
}

func (r *Router) setupSubmissionRoutes(api *gin.RouterGroup) {
	submissions := api.Group("/submissions")
	{
		submissions.GET("/:id", r.handlers.Submission.GetByID)
		submissions.POST("", r.handlers.Submission.Submit)
		submissions.PUT("/:id", r.handlers.Submission.Update)
		submissions.DELETE("/:id", r.handlers.Submission.Delete)

		// Student submissions
		api.GET("/student/submissions", r.handlers.Submission.ListMySubmissions)

		// Assignment submissions
		api.GET("/assignments/:id/submissions", r.handlers.Submission.ListByAssignment)
	}
}

func (r *Router) setupGradeRoutes(api *gin.RouterGroup) {
	grades := api.Group("/grades")
	{
		grades.GET("/:id", r.handlers.Grade.GetByID)
		grades.GET("/my", r.handlers.Grade.ListMyGrades)

		// Teacher grading
		teacherRoutes := grades.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("", r.handlers.Grade.GradeSubmission)
			teacherRoutes.PUT("/:id", r.handlers.Grade.Update)
			teacherRoutes.DELETE("/:id", r.handlers.Grade.Delete)
		}

		// Course grades
		api.GET("/courses/:id/grades", r.handlers.Grade.ListByCourse)
	}
}

func (r *Router) setupAttendanceRoutes(api *gin.RouterGroup) {
	attendance := api.Group("/attendance")
	{
		// Sessions
		attendance.GET("/sessions/:id", r.handlers.Attendance.GetSession)
		api.GET("/courses/:id/attendance/sessions", r.handlers.Attendance.ListSessionsByCourse)

		// Teacher/Admin only
		teacherRoutes := attendance.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("/sessions", r.handlers.Attendance.CreateSession)
			teacherRoutes.PUT("/sessions/:id", r.handlers.Attendance.UpdateSession)
			teacherRoutes.DELETE("/sessions/:id", r.handlers.Attendance.DeleteSession)
			teacherRoutes.POST("/sessions/:id/mark", r.handlers.Attendance.MarkAttendance)
			teacherRoutes.POST("/sessions/:id/bulk-mark", r.handlers.Attendance.BulkMarkAttendance)
		}

		// Student attendance summary (moved under /attendance to avoid conflict)
		attendance.GET("/students/:student_id/courses/:course_id/summary", r.handlers.Attendance.GetStudentCourseSummary)
	}
}

func (r *Router) setupChatRoutes(api *gin.RouterGroup) {
	chatRoutes := api.Group("/chat")
	{
		chatRoutes.GET("/rooms", r.handlers.Chat.ListMyRooms)
		chatRoutes.GET("/rooms/:id", r.handlers.Chat.GetRoom)
		chatRoutes.POST("/rooms", r.handlers.Chat.CreateRoom)
		chatRoutes.PUT("/rooms/:id", r.handlers.Chat.UpdateRoom)
		chatRoutes.DELETE("/rooms/:id", r.handlers.Chat.DeleteRoom)
		chatRoutes.POST("/rooms/direct", r.handlers.Chat.GetOrCreateDirectRoom)
		chatRoutes.GET("/rooms/:id/messages", r.handlers.Chat.ListMessages)
		chatRoutes.POST("/rooms/:id/messages", r.handlers.Chat.SendMessage)
		chatRoutes.PUT("/messages/:id", r.handlers.Chat.UpdateMessage)
		chatRoutes.DELETE("/messages/:id", r.handlers.Chat.DeleteMessage)
		chatRoutes.POST("/rooms/:id/members", r.handlers.Chat.AddMember)
		chatRoutes.DELETE("/rooms/:id/members/:user_id", r.handlers.Chat.RemoveMember)
		chatRoutes.POST("/rooms/:id/leave", r.handlers.Chat.LeaveRoom)
		chatRoutes.GET("/rooms/:id/members", r.handlers.Chat.ListMembers)
	}
}

func (r *Router) setupNotificationRoutes(api *gin.RouterGroup) {
	notifications := api.Group("/notifications")
	{
		notifications.GET("", r.handlers.Notification.ListNotifications)
		notifications.GET("/unread-count", r.handlers.Notification.GetUnreadCount)
		notifications.GET("/:id", r.handlers.Notification.GetNotification)
		notifications.POST("/:id/read", r.handlers.Notification.MarkAsRead)
		notifications.POST("/read-all", r.handlers.Notification.MarkAllAsRead)
		notifications.DELETE("/:id", r.handlers.Notification.DeleteNotification)
		notifications.DELETE("", r.handlers.Notification.DeleteAllNotifications)
	}
}

func (r *Router) setupScheduleRoutes(api *gin.RouterGroup) {
	schedule := api.Group("/schedule")
	{
		schedule.GET("", r.handlers.Schedule.ListMySchedule)
		schedule.GET("/calendar", r.handlers.Schedule.GetCalendar)
		schedule.GET("/events/:id", r.handlers.Schedule.GetByID)

		// Course schedule
		api.GET("/courses/:id/schedule", r.handlers.Schedule.ListByCourse)

		// Teacher/Admin only
		teacherRoutes := schedule.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("/events", r.handlers.Schedule.Create)
			teacherRoutes.PUT("/events/:id", r.handlers.Schedule.Update)
			teacherRoutes.DELETE("/events/:id", r.handlers.Schedule.Delete)
		}
	}
}

func (r *Router) setupAnalyticsRoutes(api *gin.RouterGroup) {
	analyticsRoutes := api.Group("/analytics")
	{
		// Track events (all users)
		analyticsRoutes.POST("/events", r.handlers.Analytics.TrackEvent)

		// Admin/Manager only
		adminRoutes := analyticsRoutes.Group("")
		adminRoutes.Use(middleware.RequireRole("admin", "manager"))
		{
			adminRoutes.GET("/system", r.handlers.Analytics.GetSystemStats)
		}

		// Course analytics (teachers)
		teacherRoutes := analyticsRoutes.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.GET("/courses/:id", r.handlers.Analytics.GetCourseStats)
			teacherRoutes.GET("/courses/:id/leaderboard", r.handlers.Analytics.GetCourseLeaderboard)
		}
	}
}

func (r *Router) setupSessionRoutes(api *gin.RouterGroup) {
	sessions := api.Group("/sessions")
	{
		sessions.GET("", r.handlers.Session.ListMySessions)
		sessions.DELETE("/:id", r.handlers.Session.RevokeSession)
		sessions.POST("/revoke-all", r.handlers.Session.RevokeAllSessions)
	}
}

func (r *Router) setupStudentRoutes(api *gin.RouterGroup) {
	students := api.Group("/students")
	{
		// Student self-management
		studentRoutes := students.Group("")
		studentRoutes.Use(middleware.RequireRole("student"))
		{
			studentRoutes.GET("/me", r.handlers.Student.GetProfile)
			studentRoutes.PUT("/me", r.handlers.Student.UpdateProfile)
			studentRoutes.GET("/me/stats", r.handlers.Student.GetStats)
		}

		// View students (teachers, managers)
		viewRoutes := students.Group("")
		viewRoutes.Use(middleware.RequireRole("teacher", "admin", "manager"))
		{
			viewRoutes.GET("", r.handlers.Student.List)
			viewRoutes.GET("/:id", r.handlers.Student.GetByID)
			viewRoutes.GET("/:id/stats", r.handlers.Student.GetStatsByID)
			viewRoutes.GET("/group/:group", r.handlers.Student.ListByGroup)
			viewRoutes.GET("/course/:id", r.handlers.Student.ListByCourse)
		}
	}
}

func (r *Router) setupTeacherRoutes(api *gin.RouterGroup) {
	teachers := api.Group("/teachers")
	{
		// Teacher self-management (allow admin too for testing)
		teacherRoutes := teachers.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.GET("/me", r.handlers.Teacher.GetProfile)
			teacherRoutes.PUT("/me", r.handlers.Teacher.UpdateProfile)
			teacherRoutes.GET("/me/stats", r.handlers.Teacher.GetStats)
			teacherRoutes.GET("/me/group-assignments", r.handlers.Group.ListTeacherAssignments) // Teacher's group assignments
		}

		// View teachers (all authenticated)
		teachers.GET("", r.handlers.Teacher.List)
		teachers.GET("/:id", r.handlers.Teacher.GetByID)

		// Admin/Manager only
		adminRoutes := teachers.Group("")
		adminRoutes.Use(middleware.RequireRole("admin", "manager"))
		{
			adminRoutes.GET("/:id/stats", r.handlers.Teacher.GetStatsByID)
			adminRoutes.GET("/:id/group-assignments", r.handlers.Group.ListTeacherAssignmentsByID) // View specific teacher's assignments
			adminRoutes.GET("/department/:department", r.handlers.Teacher.ListByDepartment)
		}
	}
}

func (r *Router) setupManagerRoutes(api *gin.RouterGroup) {
	managers := api.Group("/managers")
	managers.Use(middleware.RequireRole("manager", "admin"))
	{
		managers.GET("/me", r.handlers.Manager.GetProfile)
		managers.PUT("/me", r.handlers.Manager.UpdateProfile)
		managers.GET("/overview", r.handlers.Manager.GetSystemOverview)

		// User management
		managers.POST("/users/:id/activate", r.handlers.Manager.ActivateUser)
		managers.POST("/users/:id/deactivate", r.handlers.Manager.DeactivateUser)
		managers.POST("/users/bulk", r.handlers.Manager.BulkUserAction)

		// Admin only
		adminRoutes := managers.Group("")
		adminRoutes.Use(middleware.RequireRole("admin"))
		{
			adminRoutes.GET("", r.handlers.Manager.List)
			adminRoutes.GET("/:id", r.handlers.Manager.GetByID)
		}
	}
}

func (r *Router) setupPlagiarismRoutes(api *gin.RouterGroup) {
	plagiarismRoutes := api.Group("/plagiarism")
	plagiarismRoutes.Use(middleware.RequireRole("teacher", "admin"))
	{
		plagiarismRoutes.POST("/check", r.handlers.Plagiarism.CheckSubmission)
		plagiarismRoutes.GET("/reports/:id", r.handlers.Plagiarism.GetReport)
		plagiarismRoutes.GET("/submissions/:submission_id", r.handlers.Plagiarism.GetReportBySubmission)
		plagiarismRoutes.GET("/assignments/:assignment_id", r.handlers.Plagiarism.ListByAssignment)
		plagiarismRoutes.GET("/suspicious", r.handlers.Plagiarism.ListSuspicious)
		plagiarismRoutes.POST("/submissions/:submission_id/recheck", r.handlers.Plagiarism.Recheck)
	}
}

func (r *Router) setupUploadRoutes(api *gin.RouterGroup) {
	uploads := api.Group("/uploads")
	{
		// Public info (authenticated)
		uploads.GET("/allowed-types", r.handlers.Upload.GetAllowedTypes)

		// File operations
		uploads.POST("", r.handlers.Upload.Upload)
		uploads.POST("/url", r.handlers.Upload.UploadFromURL)
		uploads.GET("/my", r.handlers.Upload.ListMyUploads)
		uploads.GET("/:id", r.handlers.Upload.GetByID)
		uploads.DELETE("/:id", r.handlers.Upload.Delete)
		uploads.GET("/reference/:type/:reference_id", r.handlers.Upload.ListByReference)
	}
}

func (r *Router) setupGroupRoutes(api *gin.RouterGroup) {
	groups := api.Group("/groups")
	{
		// View groups (all authenticated users)
		groups.GET("", r.handlers.Group.List)
		groups.GET("/:id", r.handlers.Group.GetByID)
		groups.GET("/code/:code", r.handlers.Group.GetByCode)
		groups.GET("/:id/assignments", r.handlers.Group.ListGroupAssignments)

		// Admin/Manager only
		adminRoutes := groups.Group("")
		adminRoutes.Use(middleware.RequireRole("admin", "manager"))
		{
			adminRoutes.POST("", r.handlers.Group.Create)
			adminRoutes.PUT("/:id", r.handlers.Group.Update)
			adminRoutes.DELETE("/:id", r.handlers.Group.Delete)
			adminRoutes.POST("/assignments", r.handlers.Group.AssignTeacher)
			adminRoutes.DELETE("/assignments/:id", r.handlers.Group.UnassignTeacher)
		}
	}
}
