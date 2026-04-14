package app

import (
	"context"
	"database/sql"
	"errors"
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
	dbUrl := os.Getenv("DB_URL")
	redisAddr := os.Getenv("REDIS_ADDR")

	validateVariables(dbUrl, redisAddr)

	runMigrations(dbUrl)
	db := initDB(dbUrl)

	repo := repository.NewPostgresLinkRepository(db)
	linkCache := cache.NewRedisCache(redisAddr)
	linkService := service.NewLinkService(repo, linkCache)
	linkHandler := handler.NewLinkHandler(linkService)
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logging)
	router.Use(middleware.Recovery)

	router.Post("/links", linkHandler.CreateLink)
	router.Get("/links/{code}", linkHandler.GetByShortCode)
	router.Delete("/links/{code}", linkHandler.DeleteLinks)
	router.Get("/links", linkHandler.GetAllLinks)
	router.Get("/links/{code}/stats", linkHandler.GetStats)

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
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

func validateVariables(url string, addr string) {
	if url == "" {
		log.Fatal("DB_URL is required")
	}

	if addr == "" {
		log.Fatal("REDIS_ADDR is required")
	}
}

func initDB(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
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
