package main

import (
	"context"
	"formdata/pkg/db"
	"formdata/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()
	err := app.Bootstrap()
	if err != nil {
		log.Fatal(err)
	}
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return db.SetupSchema(app)
	})

	if err := routes.SetupRoutes(app); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := app.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ResetBootstrapState(); err != nil {
		log.Fatalf("PocketBase graceful shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped.")

}
