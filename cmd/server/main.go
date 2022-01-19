package main

import (
	"compress/flate"
	"context"
	"github.com/IlyaYP/devops/cmd/server/config"
	"github.com/IlyaYP/devops/cmd/server/handlers"
	"github.com/IlyaYP/devops/storage"
	"github.com/IlyaYP/devops/storage/infile"
	"github.com/IlyaYP/devops/storage/inmemory"
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
		cfg.StoreInterval, "STORE_FILE", cfg.StoreFile, "RESTORE", cfg.Restore)

	// Storage Q: Решил пока оставить тут, но возможно лучше перенести в config?
	var st storage.MetricStorage // Q: Это получается что? структура указатель или что?
	if cfg.StoreFile == "" {
		st = inmemory.NewMemStorage() // Q: тут я явно возврщаю указатель
	} else {
		stt, err := infile.NewFileStorage(cfg) // Q: и тут
		if err != nil {
			log.Println(err)
			return err
		}
		st = stt // Q: Что получается я созадю копию структуры или указателья???
		defer stt.Close()
	}

	// Handlers
	h := new(handlers.Handlers)
	h.St = st // Q: Тот же вопрос. Опять содается копия? (я бы не хотел полодить копии,
	// а иметь в памяти один экземпляр и передовать указатель на него)

	// Router
	r := chi.NewRouter()
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)
	r.Get("/", h.ReadHandler())
	r.Post("/update/", h.UpdateJSONHandler())
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
