package core_server

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	core_logger "RWBDwmoTask/internal/core/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux     *http.ServeMux
	config  ServerConfig
	log     *core_logger.Logger
	storage *core_domain.Storage
}

func NewServer(
	config ServerConfig,
	log *core_logger.Logger,
	storage *core_domain.Storage,
) *HTTPServer {
	return &HTTPServer{
		mux:     http.NewServeMux(),
		config:  config,
		log:     log,
		storage: storage,
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

func (s *HTTPServer) ReadEvents(
	ctx context.Context,
	js jetstream.JetStream,
) error {
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "SEARCH",
		Subjects: []string{"search.events"},
		MaxAge:   10 * time.Minute,
	})
	if err != nil {
		return err
	}

	fmt.Println("Read events stream")
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "trend-service",
		FilterSubject: "search.events",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	fmt.Println("Read events consumer")
	consumerCtx, err := consumer.Consume(func(msg jetstream.Msg) {
		var event core_domain.SearchEvent

		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			_ = msg.Ack()
			return
		}

		fmt.Println("event from nats:", event.Query)
		s.storage.Add(event.Query, 1)

		_ = msg.Ack()
	})
	if err != nil {
		return err
	}

	<-ctx.Done()
	consumerCtx.Stop()

	return err
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
