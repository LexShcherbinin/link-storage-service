package app

import (
	"context"
	"database/sql"
	"fmt"
	"link-storage-service/internal/cache"
	"link-storage-service/internal/handler"
	"link-storage-service/internal/middleware"
	"link-storage-service/internal/repository"
	"link-storage-service/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
)

type App struct {
	server *http.Server
	db     *sql.DB
	cache  cache.Cache
}

func NewApp() *App {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL is required")
	}

	fmt.Println(dbURL)

	runMigrations(dbURL)
	db := initDB()

	repo := repository.NewPostgresLinkRepository(db)

	redisAddr := os.Getenv("REDIS_ADDR")
	fmt.Println(redisAddr)
	cache := cache.NewRedisCache(redisAddr)

	service := service.NewLinkService(repo, cache)

	handler := handler.NewLinkHandler(service)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logging)
	router.Use(middleware.Recovery)

	router.Post("/links", handler.Create)
	router.Get("/links/{code}", handler.Get)
	router.Delete("/links/{code}", handler.Delete)
	router.Get("/links", handler.GetAll)
	router.Get("/links/{code}/stats", handler.Stats)

	return &App{
		server: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

func (a *App) Run() error {
	go func() {
		log.Println("server started on :8080")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if a.db != nil {
		_ = a.db.Close()
	}

	log.Println("server stopped gracefully")
	return nil
}

func initDB() *sql.DB {
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	return db
}

func runMigrations(dbUrl string) {
	m, err := migrate.New(
		"file://migrations",
		dbUrl,
	)
	if err != nil {
		log.Fatalf("migration init error: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("migration up error: %v", err)
	}

	log.Println("migrations applied")
}
