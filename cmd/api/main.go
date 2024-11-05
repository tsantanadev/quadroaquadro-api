package main

import (
	"log"

	"github.com/tsantanadev/quadroaquadro/internal/db"
	"github.com/tsantanadev/quadroaquadro/internal/env"
	filestorage "github.com/tsantanadev/quadroaquadro/internal/file_storage"
	"github.com/tsantanadev/quadroaquadro/internal/rest"
	"github.com/tsantanadev/quadroaquadro/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("PORT", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://postgres:postgres@localhost/quadroaquadro?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		TMDBConfig: TMDBConfig{
			apiKey: env.GetString("TMDB_API_KEY", ""),
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	log.Println("Database connection established")

	store := store.NewStorage(db)

	fileStorage := &filestorage.FileStorageGCP{
		Config: &filestorage.FileStorageConfig{
			Bucket: env.GetString("BUCKET", ""),
		},
	}

	app := &application{
		config:      cfg,
		store:       store,
		tmdbClient:  *rest.NewTMDBClient(cfg.TMDBConfig.apiKey),
		fileStorage: fileStorage,
	}

	mux := app.mountRoutes()

	log.Fatal(app.Run(mux))
}
