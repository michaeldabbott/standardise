package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	http.Handler
	HandleFunc(pattern string, handlerFn http.HandlerFunc)
	chi.Router
}

type Factory interface {
	Create(opts ...Option) *Server
}

type factory struct {
	config   Config
	routerFn func() Handler
}

func NewFactory(opts ...FactoryOption) Factory {
	f := &factory{
		config: defaultConfig(),
		routerFn: func() Handler {
			r := chi.NewRouter()
			return r
		},
	}

	for _, opt := range opts {
		if opt != nil {
			opt.apply(f)
		}
	}

	return f
}

type check func(http.HandlerFunc) http.HandlerFunc

type Server struct {
	Router         Handler
	healthCheck    check
	readinessCheck check
	alivenessCheck check
	config         Config
}

func (f factory) Create(opts ...Option) *Server {

	svr := &Server{
		Router: f.routerFn(),
		config: f.config,
	}

	for _, opt := range opts {
		if opt != nil {
			opt.apply(svr)
		}
	}

	svr.register()
	svr.Router.HandleFunc("/live", svr.getAlivenessHandler())
	svr.Router.HandleFunc("/ready", svr.getReadinessHandler())
	svr.Router.HandleFunc("/health", svr.getAlivenessHandler())

	return svr
}

// register registers middle on ware on the chi middleware stack
func (s *Server) register() http.Handler {
	s.Router.Use(s.timeoutMiddleware())
	s.Router.Use(s.tracingMiddleware())
	return s.Router
}

// Serve sets up a http server and starts listening
func (s *Server) Serve(ctx context.Context) error { //Take serve options
	port := s.config.Port
	if port < 1 {
		port = 8080
	}

	svr := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      s.Router,
		ReadTimeout:  time.Duration(s.config.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(s.config.WriteTimeoutMs) * time.Millisecond,
	}

	errs := make(chan error)
	go func() {
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("server failed to start up - error ", err)
			//s.logger.Error(ctx, "server failed to start up", "error", err)
			errs <- err
		} else {
			errs <- nil
		}
	}()

	log.Println("server started successfully - port ", port)

	//s.logger.Info(ctx, "server started successfully", "port", port)

	go func() {
		errs <- s.gracefulShutdown(ctx, &svr)
	}()

	return <-errs
}

// TracingMiddleware ...
func (s *Server) tracingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		//return nethttp.Middleware(s.tracer, next)
		return next
	}
}

// TimeoutMiddleware ...
func (s *Server) timeoutMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, time.Duration(s.config.RequestTimeoutSec)*time.Second, "timeout")
	}
}

//func (s *Server) getHandler(ctx context.Context) http.Handler {
//
//	s.Router.Use(s.timeoutMiddleware())
//	s.Router.Use(s.tracingMiddleware())
//	//h = s.profilingMiddleware()(h)
//	//Add other global middleware here
//	return s.Router
//}

func (s *Server) gracefulShutdown(ctx context.Context, server *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//sig := <-quit
	<-quit
	//s.logger.Info(ctx, "signal received", "signal", sig)

	timeout := time.Duration(s.config.ShutdownDelaySeconds) * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {

		//s.logger.Error(
		//	ctx,
		//	"Error while gracefully shutting down server, forcing shutdown because of error",
		//	"err", err)
		return err
	}
	//s.logger.Info(ctx, "server exited successfully")
	return nil
}

type Config struct {
	Port                 int
	ReadTimeoutMs        int
	WriteTimeoutMs       int
	RequestTimeoutSec    int
	ShutdownDelaySeconds int
	SwaggerFile          string
}

func defaultConfig() Config {
	return Config{
		Port:                 8080,
		ReadTimeoutMs:        10000,
		WriteTimeoutMs:       10000,
		RequestTimeoutSec:    10,
		ShutdownDelaySeconds: 5,
		SwaggerFile:          "/swagger.json",
	}
}
