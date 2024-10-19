package main

import (
	"log"

	"github.com/tsantanadev/social-api/internal/db"
	"github.com/tsantanadev/social-api/internal/env"
	"github.com/tsantanadev/social-api/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("PORT", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://postgres:postgres@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	log.Println("Database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mountRoutes()

	log.Fatal(app.Run(mux))
}