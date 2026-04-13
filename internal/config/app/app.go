package app

import (
	"context"
	"database/sql"
	"link-storage-service/internal/cache"
	"link-storage-service/internal/handler"
	"link-storage-service/internal/repository"
	"link-storage-service/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	server *http.Server
	db     *sql.DB
	cache  cache.Cache
}

func NewApp() *App {
	db := initDB()

	repo := repository.NewPostgresLinkRepository(db)

	cache := cache.NewRedisCache("localhost:6379")

	service := service.NewLinkService(repo, cache)

	handler := handler.NewLinkHandler(service)

	router := chi.NewRouter()

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
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:11432/links?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	return db
}
