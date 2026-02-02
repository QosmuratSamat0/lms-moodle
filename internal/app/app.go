package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	deliveryhttp "ap1-final-mini-moodle/internal/delivery/http"
	"ap1-final-mini-moodle/internal/repository"
	"ap1-final-mini-moodle/internal/usecase"
)

func Run() error {
	store := repository.NewInMemoryStore()
	teacherUsecase := usecase.NewTeacherUsecase(store)
	adminUsecase := usecase.NewAdminUsecase(store)

	router := deliveryhttp.NewRouter(teacherUsecase, adminUsecase)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	errs := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errs:
		return err
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
}
