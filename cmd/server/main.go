package main

import (
	"context"
	"flag"
	"github.com/IlyaYP/devops/cmd/server/config"
	"github.com/IlyaYP/devops/cmd/server/handlers"
	"github.com/IlyaYP/devops/storage"
	"github.com/IlyaYP/devops/storage/infile"
	"github.com/IlyaYP/devops/storage/inmemory"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cfg config.Config

func init() {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Server address")
	//flag.IntVar(&cfg.StoreInterval, "i", 300, "Store interval in seconds")
	flag.DurationVar(&cfg.StoreInterval, "i", time.Duration(300)*time.Second, "Store interval in seconds")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "Store file")
	flag.BoolVar(&cfg.Restore, "r", true, "Restore data from file when start")
}

func main() {
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("Server start using args:ADDRESS", cfg.Address, "STORE_INTERVAL",
		cfg.StoreInterval, "STORE_FILE", cfg.StoreFile, "RESTORE", cfg.Restore)

	var st storage.MetricStorage // Q: Это получается что? структура указатель или что?
	if cfg.StoreFile == "" {
		st = inmemory.NewMemStorage() // Q: тут я явно возврщаю указатель
	} else {
		stt, err := infile.NewFileStorage(&cfg) // Q: и тут
		if err != nil {
			log.Fatal(err)
		}
		st = stt // Q: Что получается я созадю копию структуры или указателья???
		defer stt.Close()
	}

	r := chi.NewRouter()
	h := new(handlers.Handlers)
	h.St = st // Q: Тот же вопрос. Опять содается копия? (я бы не хотел полодить копии,
	// а иметь в памяти один экземпляр и передовать указатель на него)
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
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	log.Println("timeout of 5 seconds.")
	<-ctx.Done()

	log.Println("Server exiting")

}
