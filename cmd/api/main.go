package main

import (
	"context"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JasonLeonnn/paygard/internal/config"
	"github.com/JasonLeonnn/paygard/internal/db"
	"github.com/JasonLeonnn/paygard/internal/http"
	"github.com/JasonLeonnn/paygard/internal/http/middleware"
	"github.com/JasonLeonnn/paygard/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()

	// Load .env if present (ignored if missing)
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	transactionService := services.NewTransactionService(pool)
	alertService := services.NewAlertService(pool)

	r := chi.NewRouter()
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.MetricsMiddleware)

	r.Get("/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Handle("/metrics", promhttp.Handler())

	r.Post("/transactions", http.CreateTransactionHandler(transactionService))

	r.Get("/alerts", http.GetAlertsHandler(alertService))

	server := &nethttp.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  150 * time.Second,
	}

	// Periodic baselines updater
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := db.UpdateBaselines(context.Background(), pool, cfg.BaselineWindowDays); err != nil {
					log.Printf("Failed to update baselines: %v", err)
				}
			case <-stop:
				return
			}
		}
	}()

	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// stop background routines
	close(stop)

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server exited properly")
}
