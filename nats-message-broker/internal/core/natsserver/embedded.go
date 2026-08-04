package natsserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

const MonitorPath = "/nats"

type EmbeddedServer struct {
	server *server.Server
	host   string
	port   int
}

func Start(host string, port int, httpPort int) (*EmbeddedServer, error) {
	// NATS uses -1 for auto-assign; port 0 is not supported and never becomes ready.
	if port == 0 {
		port = -1
	}
	if httpPort == 0 {
		httpPort = -1
	}

	opts := &server.Options{
		Host:     host,
		Port:     port,
		HTTPHost: host,
		HTTPPort: httpPort,
		NoSigs:   true,
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("create nats server: %w", err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) {
		ns.Shutdown()
		if port > 0 {
			return nil, fmt.Errorf(
				"embedded nats server not ready on %s:%d: port may already be in use; stop the other process or set NATS_PORT=-1 for auto-assign",
				host, port,
			)
		}
		return nil, fmt.Errorf("embedded nats server not ready on %s (auto-assign failed)", host)
	}

	actualPort := resolvePort(ns, port)
	if actualPort <= 0 {
		ns.Shutdown()
		return nil, fmt.Errorf("embedded nats server started but bound port is unknown")
	}

	log.Printf("Embedded NATS server started on %s:%d", host, actualPort)

	return &EmbeddedServer{
		server: ns,
		host:   host,
		port:   actualPort,
	}, nil
}

func resolvePort(ns *server.Server, configuredPort int) int {
	if addr := ns.Addr(); addr != nil {
		if tcpAddr, ok := addr.(*net.TCPAddr); ok && tcpAddr.Port > 0 {
			return tcpAddr.Port
		}
	}
	if configuredPort > 0 {
		return configuredPort
	}
	return 0
}

func (s *EmbeddedServer) Port() int {
	return s.port
}

func (s *EmbeddedServer) IsRunning() bool {
	return s.server != nil && s.server.Running()
}

func (s *EmbeddedServer) HTTPHandler() http.Handler {
	if s.server == nil {
		return nil
	}
	return s.server.HTTPHandler()
}

func (s *EmbeddedServer) Shutdown() {
	if s.server != nil {
		s.server.Shutdown()
		log.Println("Embedded NATS server stopped")
	}
}

func (s *EmbeddedServer) ClientURL() string {
	return fmt.Sprintf("nats://%s:%d", s.host, s.port)
}
