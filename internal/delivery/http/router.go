package deliveryhttp

import (
	"ap1-final-mini-moodle/internal/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(teacher *usecase.TeacherUsecase, admin *usecase.AdminUsecase) *gin.Engine {
	handler := NewHandler(teacher, admin)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/health", handler.Health)

	teacherGroup := engine.Group("/api/teacher", Auth(), RequireRole("teacher"))
	teacherGroup.POST("/courses", handler.CreateCourse)
	teacherGroup.POST("/courses/:courseId/assignments", handler.CreateAssignment)
	teacherGroup.GET("/assignments/:assignmentId/submissions", handler.ListSubmissions)
	teacherGroup.POST("/submissions/:submissionId/grade", handler.GradeSubmission)

	adminGroup := engine.Group("/api/admin", Auth(), RequireRole("admin"))
	adminGroup.GET("/courses", handler.AdminCourses)

	return engine
}
