package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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

	httpdelivery "github.com/ap1-final-mini-moodle/internal/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Initialize database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/moodle"
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	// Initialize repositories
	userRepository := userRepo.NewPostgresRepository(dbPool)
	courseRepository := courseRepo.NewPostgresRepository(dbPool)
	enrollmentRepository := enrollmentRepo.NewPostgresRepository(dbPool)
	assignmentRepository := assignmentRepo.NewPostgresRepository(dbPool)
	submissionRepository := submissionRepo.NewPostgresRepository(dbPool)
	gradeRepository := gradeRepo.NewPostgresRepository(dbPool)
	attendanceRepository := attendanceRepo.NewPostgresRepository(dbPool)
	notificationRepository := notificationRepo.NewPostgresRepository(dbPool)
	chatRepository := chatRepo.NewPostgresRepository(dbPool)
	uploadRepository := uploadRepo.NewPostgresRepository(dbPool)

	// Initialize use cases
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

	// Initialize handlers
	userHandler := httpdelivery.NewUserHandler(userService)
	courseHandler := httpdelivery.NewCourseHandler(courseService)
	enrollmentHandler := httpdelivery.NewEnrollmentHandler(enrollmentService)
	assignmentHandler := httpdelivery.NewAssignmentHandler(assignmentService)
	submissionHandler := httpdelivery.NewSubmissionHandler(submissionService)
	gradeHandler := httpdelivery.NewGradeHandler(gradeService)
	attendanceHandler := httpdelivery.NewAttendanceHandler(attendanceService)
	notificationHandler := httpdelivery.NewNotificationHandler(notificationService)
	chatHandler := httpdelivery.NewChatHandler(chatService)
	uploadHandler := httpdelivery.NewUploadHandler(uploadService)

	// Setup router
	router := gin.Default()

	// Routes
	api := router.Group("/api/v1")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("", userHandler.Register)
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.GetByID)
			users.PATCH("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

		// Course routes
		courses := api.Group("/courses")
		{
			courses.POST("", courseHandler.Create)
			courses.GET("", courseHandler.List)
			courses.GET("/:id", courseHandler.GetByID)
			courses.PATCH("/:id", courseHandler.Update)
			courses.DELETE("/:id", courseHandler.Delete)
		}

		// Enrollment routes
		enrollments := api.Group("/enrollments")
		{
			enrollments.POST("", enrollmentHandler.Enroll)
			enrollments.GET("/course/:courseID", enrollmentHandler.ListByCourse)
			enrollments.GET("/student/:studentID", enrollmentHandler.ListByStudent)
			enrollments.DELETE("/:id", enrollmentHandler.Remove)
		}

		// Assignment routes
		assignments := api.Group("/assignments")
		{
			assignments.POST("", assignmentHandler.Create)
			assignments.GET("/:id", assignmentHandler.GetByID)
			assignments.GET("/course/:courseID", assignmentHandler.ListByCourse)
			assignments.DELETE("/:id", assignmentHandler.Delete)
		}

		// Submission routes
		submissions := api.Group("/submissions")
		{
			submissions.POST("", submissionHandler.Submit)
			submissions.GET("/:id", submissionHandler.GetByID)
			submissions.GET("/assignment/:assignmentID", submissionHandler.ListByAssignment)
			submissions.GET("/student/:studentID", submissionHandler.ListByStudent)
			submissions.DELETE("/:id", submissionHandler.Delete)
		}

		// Grade routes
		grades := api.Group("/grades")
		{
			grades.POST("", gradeHandler.Grade)
			grades.GET("/:id", gradeHandler.GetByID)
			grades.DELETE("/:id", gradeHandler.Delete)
		}

		// Attendance routes
		attendance := api.Group("/attendance")
		{
			attendance.POST("", attendanceHandler.Record)
			attendance.GET("/course/:courseID", attendanceHandler.ListByCourse)
			attendance.DELETE("/:id", attendanceHandler.Delete)
		}

		// Notification routes
		notifications := api.Group("/notifications")
		{
			notifications.GET("/user/:userID", notificationHandler.ListByUser)
			notifications.PATCH("/:id/read", notificationHandler.MarkAsRead)
			notifications.DELETE("/:id", notificationHandler.Delete)
		}

		// Chat routes
		chat := api.Group("/chat")
		{
			chat.POST("", chatHandler.SendMessage)
			chat.GET("/course/:courseID", chatHandler.ListByCourse)
			chat.DELETE("/:id", chatHandler.Delete)
		}

		// Upload routes
		uploads := api.Group("/uploads")
		{
			uploads.POST("", uploadHandler.Upload)
			uploads.GET("/:id", uploadHandler.GetByID)
			uploads.GET("/user/:userID", uploadHandler.ListByUser)
			uploads.DELETE("/:id", uploadHandler.Delete)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Server error:", err)
	}
}
