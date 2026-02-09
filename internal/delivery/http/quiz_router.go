package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type QuizModule struct {
	handler     *QuizHandler
	authService *authUC.Service
}

func NewQuizModule(handler *QuizHandler, authService *authUC.Service) *QuizModule {
	return &QuizModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *QuizModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthTokenMiddleware(m.authService))

	quizzes := api.Group("/quizzes")
	{
		quizzes.POST("", m.handler.Create)
		quizzes.PUT("/:id", m.handler.Update)
		quizzes.DELETE("/:id", m.handler.Delete)
		quizzes.GET("/:id", m.handler.GetByID)
		quizzes.GET("/:id/questions", m.handler.GetQuizQuestions)
		quizzes.POST("/:id/attempts", m.handler.StartAttempt)
		quizzes.GET("/:id/attempts", m.handler.GetStudentAttempts)
	}

	attempts := api.Group("/quiz-attempts")
	{
		attempts.POST("/:attemptID/submit", m.handler.SubmitAttempt)
		attempts.GET("/:attemptID/result", m.handler.GetAttemptResult)
	}

	api.GET("/courses/:courseID/quizzes", m.handler.ListByCourse)
}
