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
	router := handler.NewRouter(linkHandler)

	return &App{
		server: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
		db:    db,
		cache: linkCache,
	}
}

func (a *App) Start() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.db.Close(); err != nil {
		log.Println("error closing db:", err)
	}

	if rc, ok := a.cache.(interface{ Close() error }); ok {
		if err := rc.Close(); err != nil {
			log.Println("error closing cache:", err)
		}
	}

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
