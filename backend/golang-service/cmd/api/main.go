// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/api"
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
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/config"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/database"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/websocket"
	cloudupload "github.com/MaqsattoTeam/aLMS/golang-service/pkg/upload"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	log.Printf("Starting %s in %s mode", cfg.App.AppName, cfg.App.Environment)

	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize PostgreSQL
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("✓ Connected to PostgreSQL")

	// Initialize Redis (optional, handle gracefully if not available)
	redisClient, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Continuing without Redis...")
	} else {
		defer redisClient.Close()
		log.Println("✓ Connected to Redis")
	}

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT)

	// Initialize repositories
	userRepo := user.NewRepository(db.Pool)
	courseRepo := course.NewRepository(db.Pool)
	enrollmentRepo := enrollment.NewRepository(db.Pool)
	assignmentRepo := assignment.NewRepository(db.Pool)
	submissionRepo := submission.NewRepository(db.Pool)
	gradeRepo := grade.NewRepository(db.Pool)
	attendanceRepo := attendance.NewRepository(db.Pool)
	chatRepo := chat.NewRepository(db.Pool)
	notificationRepo := notification.NewRepository(db.Pool)
	scheduleRepo := schedule.NewRepository(db.Pool)
	analyticsRepo := analytics.NewRepository(db.Pool)
	sessionRepo := session.NewRepository(db.Pool)
	studentRepo := student.NewRepository(db.Pool)
	teacherRepo := teacher.NewRepository(db.Pool)
	managerRepo := manager.NewRepository(db.Pool)
	plagiarismRepo := plagiarism.NewRepository(db.Pool)
	uploadRepo := upload.NewRepository(db.Pool)
	groupRepo := group.NewRepository(db.Pool)

	// Initialize plagiarism detector
	plagiarismDetector := plagiarism.NewDetector(5)

	// Initialize Cloudinary uploader
	var cloudinaryUploader *cloudupload.CloudinaryUploader
	if cfg.Cloudinary.CloudName != "" && cfg.Cloudinary.APIKey != "" && cfg.Cloudinary.APISecret != "" {
		var err error
		cloudinaryUploader, err = cloudupload.NewCloudinaryUploader(cloudupload.Config{
			CloudName:    cfg.Cloudinary.CloudName,
			APIKey:       cfg.Cloudinary.APIKey,
			APISecret:    cfg.Cloudinary.APISecret,
			UploadPreset: cfg.Cloudinary.UploadPreset,
			Folder:       cfg.Cloudinary.Folder,
			MaxFileSize:  cfg.Cloudinary.MaxFileSize,
		})
		if err != nil {
			log.Printf("Warning: Failed to initialize Cloudinary: %v", err)
		} else {
			log.Println("✓ Cloudinary initialized")
		}
	} else {
		log.Println("Warning: Cloudinary not configured, file uploads will be disabled")
	}

	// Initialize services
	userService := user.NewService(userRepo, jwtManager)
	courseService := course.NewService(courseRepo)
	enrollmentService := enrollment.NewService(enrollmentRepo)
	assignmentService := assignment.NewService(assignmentRepo)
	submissionService := submission.NewService(submissionRepo)
	gradeService := grade.NewService(gradeRepo)
	attendanceService := attendance.NewService(attendanceRepo)
	chatService := chat.NewService(chatRepo)
	notificationService := notification.NewService(notificationRepo)
	scheduleService := schedule.NewService(scheduleRepo)
	analyticsService := analytics.NewService(analyticsRepo)
	sessionService := session.NewService(sessionRepo)
	studentService := student.NewService(studentRepo)
	teacherService := teacher.NewService(teacherRepo)
	managerService := manager.NewService(managerRepo)
	plagiarismService := plagiarism.NewService(plagiarismRepo, plagiarismDetector)
	uploadService := upload.NewService(uploadRepo, cloudinaryUploader)
	groupService := group.NewService(groupRepo)

	// Initialize handlers
	handlers := &api.Handlers{
		User:         user.NewHandler(userService),
		Course:       course.NewHandler(courseService),
		Enrollment:   enrollment.NewHandler(enrollmentService),
		Assignment:   assignment.NewHandler(assignmentService),
		Submission:   submission.NewHandler(submissionService),
		Grade:        grade.NewHandler(gradeService),
		Attendance:   attendance.NewHandler(attendanceService),
		Chat:         chat.NewHandler(chatService),
		Notification: notification.NewHandler(notificationService),
		Schedule:     schedule.NewHandler(scheduleService),
		Analytics:    analytics.NewHandler(analyticsService),
		Session:      session.NewHandler(sessionService),
		Student:      student.NewHandler(studentService),
		Teacher:      teacher.NewHandler(teacherService),
		Manager:      manager.NewHandler(managerService),
		Plagiarism:   plagiarism.NewHandler(plagiarismService),
		Upload:       upload.NewHandler(uploadService),
		Group:        group.NewHandler(groupService),
	}

	// Initialize Gin router
	engine := gin.New()
	api.SetupMiddleware(engine)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Setup router with all handlers
	router := api.NewRouter(engine, authMiddleware, handlers)
	router.SetupRoutes()

	// Initialize WebSocket hub and start it
	hub := websocket.NewHub()
	go hub.Run()

	// Setup WebSocket routes for chat
	routes.SetupChatWebSocket(engine, authMiddleware, chatService, hub)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
