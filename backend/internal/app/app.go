package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ap1-final-mini-moodle/internal/shared/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDatabase() *pgxpool.Pool {
	cfg := config.Load()
	return InitDatabaseWithConfig(cfg)
}

func InitDatabaseWithConfig(cfg *config.Config) *pgxpool.Pool {
	ctx := context.Background()
	dbURL := cfg.GetDatabaseURL()

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
	return pool
}

type App struct {
	db             *pgxpool.Pool
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	notificationCh chan Notification
	backgroundDone chan struct{}
}

type Notification struct {
	UserID    string
	Message   string
	Type      string
	Timestamp time.Time
}

func New(db *pgxpool.Pool) *App {
	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		db:             db,
		ctx:            ctx,
		cancel:         cancel,
		notificationCh: make(chan Notification, 100),
		backgroundDone: make(chan struct{}),
	}

	app.startBackgroundWorkers()

	return app
}

func (a *App) startBackgroundWorkers() {
	a.wg.Add(1)
	go a.notificationWorker()

	a.wg.Add(1)
	go a.healthCheckWorker()
}

func (a *App) notificationWorker() {
	defer a.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			log.Println("Notification worker shutting down...")
			return

		case notification := <-a.notificationCh:
			a.processNotification(notification)

		case <-ticker.C:
			a.processPendingNotifications()
		}
	}
}

func (a *App) healthCheckWorker() {
	defer a.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			log.Println("Health check worker shutting down...")
			return

		case <-ticker.C:
			a.checkDatabaseHealth()
		}
	}
}

func (a *App) processNotification(n Notification) {
	log.Printf("[NOTIFICATION] User: %s, Type: %s, Message: %s", n.UserID, n.Type, n.Message)
}

func (a *App) processPendingNotifications() {
	log.Println("[BACKGROUND] Processing pending notifications...")
}

func (a *App) checkDatabaseHealth() {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	err := a.db.Ping(ctx)
	if err != nil {
		log.Printf("[ERROR] Database health check failed: %v", err)
	} else {
		log.Println("[HEALTH] Database connection OK")
	}
}

func (a *App) SendNotification(userID, message, notificationType string) {
	notification := Notification{
		UserID:    userID,
		Message:   message,
		Type:      notificationType,
		Timestamp: time.Now(),
	}

	select {
	case a.notificationCh <- notification:
	case <-a.ctx.Done():
		log.Println("App context cancelled, notification not sent")
	default:
		log.Println("Notification channel full, dropping notification")
	}
}

func (a *App) Close() error {
	log.Println("Closing app...")
	a.cancel()

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Background workers stopped")
	case <-time.After(10 * time.Second):
		log.Println("Timeout waiting for background workers to stop")
	}

	a.db.Close()
	return nil
}
