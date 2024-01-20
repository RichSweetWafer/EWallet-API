package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RichSweetWafer/EWallet-API/config"
	"github.com/RichSweetWafer/EWallet-API/wallets"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg     config.HTTPServer
	wallets wallets.Interface
	router  *chi.Mux
}

func NewServer(cfg config.HTTPServer, wallets wallets.Interface) *Server {
	server := Server{
		cfg:     cfg,
		wallets: wallets,
		router:  chi.NewRouter(),
	}
	server.router.Use(middleware.Logger)
	server.routes()
	log.Println("Chi server started.")
	return &server
}

func (s *Server) Start(ctx context.Context) {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Port),
		Handler: s.router,
	}

	shutdownComplete := handleShutdown(func() {
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server.Shutdown failed: %v\n", err)
		}
	})

	if err := server.ListenAndServe(); err == http.ErrServerClosed {
		<-shutdownComplete
	} else {
		log.Printf("http.ListenAndServe failed: %v\n", err)
	}

	log.Println("Shutdown complete")
}

func handleShutdown(onShutDownSignal func()) <-chan struct{} {
	shutdown := make(chan struct{})

	go func() {
		shutDownSignal := make(chan os.Signal, 1)
		signal.Notify(shutDownSignal, os.Interrupt, syscall.SIGTERM)

		<-shutDownSignal

		onShutDownSignal()

		close(shutdown)
	}()

	return shutdown
}
