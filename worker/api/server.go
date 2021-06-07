package api

import (
	"fmt"
	"net"

	"google.golang.org/grpc/credentials"
)

func newServer(cfg Config, cred credentials.TransportCredentials) *Server {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port))
}
