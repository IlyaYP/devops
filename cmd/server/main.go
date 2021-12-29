package main

import (
	"context"
	"github.com/IlyaYP/devops/cmd/server/handlers"
	"github.com/IlyaYP/devops/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	st := inmemory.NewStorage()
	//testStore(st)

	r := chi.NewRouter()
	r.Get("/", handlers.ReadHandler(st))
	r.Get("/value/{MType}/{MName}", handlers.GetHandler(st))
	r.Post("/update/{MType}/{MName}/{MVal}", handlers.UpdateHandler(st))

	srv := &http.Server{
		Addr:    ":8080",
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
