package core_server

import (
	core_logger "RWBDwmoTask/internal/core/logger"
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config ServerConfig
	log    *core_logger.Logger
}

func NewServer(
	config ServerConfig,
	log *core_logger.Logger,
) *HTTPServer {
	return &HTTPServer{
		mux:    http.NewServeMux(),
		config: config,
		log:    log,
	}
}

func (s *HTTPServer) ResisterApiVersionRouter(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.version)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}

}

func (s *HTTPServer) AddFrond() {
	s.mux.Handle(
		"GET /css/",
		http.StripPrefix(
			"/css/",
			http.FileServer(http.Dir("./public/css")),
		),
	)

	s.mux.Handle(
		"GET /js/",
		http.StripPrefix(
			"/js/",
			http.FileServer(http.Dir("./public/js")),
		),
	)
}

func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := route.Method + " " + route.Path
		s.mux.Handle(pattern, route.Handler)
	}
}

func (s *HTTPServer) Start(ctx context.Context) error {

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.mux,
	}

	ch := make(chan error)

	go func() {
		defer close(ch)

		s.log.Warn("Starting HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server error: %w", err)
		}

	case <-ctx.Done():
		s.log.Warn("Stopping HTTP server", zap.Error(ctx.Err()))

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown server error: %w", err)
		}
		s.log.Warn("HTTP shutdown server stopped")
	}

	return nil
}
