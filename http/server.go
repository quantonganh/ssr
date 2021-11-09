package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/quantonganh/ssr"
)

const (
	shutdownTimeout = 3 * time.Second
)

type Server struct {
	ln net.Listener
	Addr string

	server *http.Server
	router *mux.Router

	RepositoryService ssr.RepositoryService
	ScanService ssr.ScanService
}

func NewServer(repositoryService ssr.RepositoryService, scanService ssr.ScanService) *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),

		RepositoryService: repositoryService,
		ScanService: scanService,
	}

	zlog := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
	s.router.Use(hlog.NewHandler(zlog))
	s.router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	s.router.Use(hlog.RemoteAddrHandler("ip"))
	s.router.Use(hlog.UserAgentHandler("user_agent"))
	s.router.Use(hlog.RefererHandler("referer"))
	s.router.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	s.router.HandleFunc("/scans/{repoID}", s.Error(s.CreateScanHandler)).Methods(http.MethodPost)
	s.router.HandleFunc("/scans/{scanID}", s.Error(s.GetScanHandler)).Methods(http.MethodGet)
	s.router.HandleFunc("/scans/{scanID}", s.Error(s.UpdateScanHandler)).Methods(http.MethodPut)
	s.router.HandleFunc("/scans/{scanID}", s.Error(s.DeleteScanHandler)).Methods(http.MethodDelete)
	s.router.HandleFunc("/scans", s.Error(s.ListScansHandler)).Methods(http.MethodGet)

	return s
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Open() (err error) {
	s.ln, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return errors.Errorf("failed to listen to port %s: %v", s.Addr, err)
	}

	go func() {
		_ = s.server.Serve(s.ln)
	}()

	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}