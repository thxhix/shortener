package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/interfaces"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServerInterface the contract for running the HTTP server.
type ServerInterface interface {
	// StartPooling starts the HTTP server and begins listening for requests.
	// If profiling is enabled in the config, a separate pprof server is also started.
	StartPooling() error

	// StopServer gracefully shuts down the running HTTP server.
	// It stops accepting new connections and waits for in-flight requests
	// to finish until the provided timeout expires. If the timeout is reached
	// before all requests are completed, the server is forcefully closed.
	// The method returns any error encountered during shutdown.
	StopServer(ctx context.Context) error

	startWithHTTP() error
	startWithHTTPS() error
}

// Server is the concrete implementation of the HTTP server.
// It holds the configuration, router, database and logger instances.
type Server struct {
	http     *http.Server
	config   *config.Config
	router   *chi.Mux
	database interfaces.Database
	logger   *zap.SugaredLogger
}

// NewServer creates a new Server instance with the provided configuration,
// router, database, and logger. It returns a ServerInterface implementation.
func NewServer(config config.Config, router chi.Mux, db interfaces.Database, logger *zap.SugaredLogger) ServerInterface {
	server := http.Server{
		Addr:    config.Address,
		Handler: &router,
	}

	return &Server{
		http:     &server,
		config:   &config,
		router:   &router,
		database: db,
		logger:   logger,
	}
}

// StartPooling starts the main HTTP API server using the configured address.
// If profiling is enabled in the configuration, a separate pprof server is
// started on the ProfilerAddress in a separate goroutine. The method blocks
// until the main HTTP server exits or encounters an error.
func (s *Server) StartPooling() error {
	s.logger.Info("* * * Запускаюсь * * *")
	s.logger.Infof("Адрес: %s", s.config.Address)
	s.logger.Infof("Base URL: %s", s.config.BaseURL)
	if s.config.EnableProfiler {
		s.logger.Infof("ProfilerAddress: %s", s.config.ProfilerAddress)
	}
	s.logger.Info("* * * * * * * * * * *")

	if s.config.EnableProfiler {
		// pprof отдельно, чтобы не мешать в API
		go func() {
			addr := s.config.ProfilerAddress
			if err := http.ListenAndServe(addr, nil); err != nil {
				s.logger.Errorf("pprof server error: %v", err)
			}
		}()
	}

	idleConnsClosed := make(chan struct{})

	sigint := make(chan os.Signal, 3)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.StopServer(ctx); err != nil {
			// server stop error
			s.logger.Errorf("HTTP server Shutdown: %v", err)
		}

		err := s.database.Close()
		if err != nil {
			s.logger.Errorf("database Close: %v", err)
		}

		close(idleConnsClosed)
	}()

	if !s.config.EnableHTTPS {
		err := s.startWithHTTP()
		if err != nil {
			return err
		}
	} else {
		err := s.startWithHTTPS()
		if err != nil {
			return err
		}
	}

	<-idleConnsClosed

	return nil
}

func (s *Server) startWithHTTP() error {
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) startWithHTTPS() error {
	cert, err := getTLSCert()
	if err != nil {
		return err
	}

	cfg := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}

	ln, err := tls.Listen("tcp", s.config.Address, cfg)
	if err != nil {
		return fmt.Errorf("tls listen: %w", err)
	}

	s.logger.Infof("HTTPS server started on %s", s.config.Address)

	if err := s.http.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// StopServer gracefully shuts down the running HTTP server.
// It stops accepting new connections and waits for in-flight requests
// to finish until the provided timeout expires. If the timeout is reached
// before all requests are completed, the server is forcefully closed.
// The method returns any error encountered during shutdown.
func (s *Server) StopServer(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
