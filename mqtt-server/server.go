package mqttserver

import (
	"space-game_mk4/mqtt-server/hooks"
	"log"
	"log/slog"
	"os"

	"github.com/cockroachdb/pebble"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	ph "github.com/mochi-mqtt/server/v2/hooks/storage/pebble"
	"github.com/mochi-mqtt/server/v2/listeners"
)

type Server struct {
	server *mqtt.Server
	db     *pebble.DB
}

func (s *Server) Close() {
	s.server.Close()
}

func (s *Server) Serve() error {
	return s.server.Serve()
}

func NewServer(opts ...func(*Server)) *Server {
	s := &Server{
		server: mqtt.New(&mqtt.Options{
			InlineClient: true,
		}),
	}
	s.server.Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: nil, //nil == info
	}))

	for _, o := range opts {
		o(s)
	}
	return s
}

func WithPebbleDB(db *pebble.DB) func(*Server) {
	return func(s *Server) {
		s.db = db
	}
}

func WithPebbleHook(path string) func(*Server) {
	return func(s *Server) {
		err := s.server.AddHook(new(ph.Hook), &ph.Options{
			Path: path,
			Mode: ph.NoSync,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WithOpenAuthHook() func(*Server) {
	return func(s *Server) {
		if err := s.server.AddHook(new(auth.AllowHook), nil); err != nil {
			log.Fatal(err)
		}
	}
}

func WithSystemHook() func(*Server) {
	return func(s *Server) {
		if err := s.server.AddHook(new(hooks.SystemHook), &hooks.SystemHooksOptions{
			Server: s.server,
			DB:     s.db,
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func WithGameHook() func(*Server) {
	return func(s *Server) {
		if err := s.server.AddHook(new(hooks.GameHook), &hooks.GameHookOptions{
			Server: s.server,
			DB:     s.db,
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func WithUserHook() func(*Server) {
	return func(s *Server) {
		if err := s.server.AddHook(new(hooks.UserHook), &hooks.UserHookOptions{
			Server: s.server,
			DB:     s.db,
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func WithMarketHook() func(*Server) {
	return func(s *Server) {
		if err := s.server.AddHook(new(hooks.MarketHook), &hooks.MarketHookOptions{
			Server: s.server,
			DB:     s.db,
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func WithDefaultListener() func(*Server) {
	return func(s *Server) {
		tcp := listeners.NewTCP(listeners.Config{
			ID:      "server",
			Address: ":1883",
		})
		err := s.server.AddListener(tcp)
		if err != nil {
			log.Fatal(err)
		}

	}
}
