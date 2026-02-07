package http

import (
	"github.com/gin-gonic/gin"
)

type QuizModule struct {
	handler *QuizHandler
	secret  []byte
}

func NewQuizModule(handler *QuizHandler, secret []byte) *QuizModule {
	return &QuizModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *QuizModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	quizzes := api.Group("/quizzes")
	{
		// Teacher routes
		quizzes.POST("", m.handler.Create)
		quizzes.PUT("/:id", m.handler.Update)
		quizzes.DELETE("/:id", m.handler.Delete)

		// Common routes
		quizzes.GET("/:id", m.handler.GetByID)
		quizzes.GET("/:id/questions", m.handler.GetQuizQuestions)

		// Student routes
		quizzes.POST("/:id/attempts", m.handler.StartAttempt)
		quizzes.GET("/:id/attempts", m.handler.GetStudentAttempts)
	}

	// Attempt routes
	attempts := api.Group("/quiz-attempts")
	{
		attempts.POST("/:attemptID/submit", m.handler.SubmitAttempt)
		attempts.GET("/:attemptID/result", m.handler.GetAttemptResult)
	}

	// Course-specific quiz routes
	api.GET("/courses/:courseID/quizzes", m.handler.ListByCourse)
}
