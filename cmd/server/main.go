package main

import (
	"compress/flate"
	"context"
	"github.com/IlyaYP/devops/cmd/server/config"
	"github.com/IlyaYP/devops/cmd/server/handlers"
	"github.com/IlyaYP/devops/storage"
	"github.com/IlyaYP/devops/storage/infile"
	"github.com/IlyaYP/devops/storage/inmemory"
	"github.com/IlyaYP/devops/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Server start using args:ADDRESS", cfg.Address, "STORE_INTERVAL",
		cfg.StoreInterval, "STORE_FILE", cfg.StoreFile, "RESTORE", cfg.Restore, "KEY",
		cfg.Key, "DATABASE_DSN", cfg.DBDsn)

	var st storage.MetricStorage
	if cfg.DBDsn == "" {
		if cfg.StoreFile == "" {
			st = inmemory.NewMemStorage()
		} else {
			stt, err := infile.NewFileStorage(context.Background(), cfg)
			if err != nil {
				log.Println(err)
				return err
			}
			st = stt
			defer stt.Close(context.Background())
		}
	} else {
		stt, err := postgres.NewPostgres(context.Background(), cfg.DBDsn)
		if err != nil {
			log.Println(err)
			return err
		}
		st = stt
		defer stt.Close()
	}
	// Handlers
	h := new(handlers.Handlers)
	h.Key = cfg.Key
	h.St = st

	// Router
	r := chi.NewRouter()
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)
	r.Get("/", h.ReadHandler())
	r.Get("/ping", h.Ping())
	r.Post("/update/", h.UpdateJSONHandler())
	r.Post("/updates/", h.UpdatesJSONHandler())
	r.Post("/value/", h.GetJSONHandler())
	r.Get("/value/{MType}/{MName}", h.GetHandler())
	r.Post("/update/{MType}/{MName}/{MVal}", h.UpdateHandler())

	srv := &http.Server{
		Addr:    cfg.Address, //":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
		return err
	}

	log.Println("timeout of 5 seconds.")
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()

	log.Println("Server exiting")
	return nil
}
