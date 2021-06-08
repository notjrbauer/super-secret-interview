package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/notjrbauer/interview/super-secret-interview/proto"
	"github.com/notjrbauer/interview/super-secret-interview/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	ls  net.Listener
	cli *grpc.Server
}

func loadTLSCredentials(cfg *worker.Config) (credentials.TransportCredentials, error) {
	// load CA cert
	pemClientCA, err := ioutil.ReadFile(cfg.Server.CACert)
	if err != nil {
		return nil, err
	}
	// create pool to add client CA's certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	// load server cert and key
	serverCert, err := tls.LoadX509KeyPair(cfg.Server.SSLCert, cfg.Server.SSLKey)
	if err != nil {
		return nil, err
	}
	// configure cert termination.
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13, // ***
	}
	return credentials.NewTLS(config), nil
}

func newServer(cfg *worker.Config, cred credentials.TransportCredentials) (*server, error) {
	ls, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port))
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(cred),
		grpc.UnaryInterceptor(AuthzInterceptor),
	)
	proto.RegisterWorkerServiceServer(grpcServer, &workerServer{
		Service: worker.NewService(cfg),
	})

	s := &server{}
	s.ls = ls
	s.cli = grpcServer
	return s, nil
}

// ListenAndServer takes a worker config and serves the api.
func ListenAndServe(cfg *worker.Config) error {
	creds, err := loadTLSCredentials(cfg)
	if err != nil {
		return fmt.Errorf("error loading tls creds: %w", err)
	}

	s, err := newServer(cfg, creds)
	if err != nil {
		return fmt.Errorf("error loading configurating server: %w", err)
	}

	defer s.Close()
	log.Printf("now serving @ %s:%d\n", cfg.Server.Hostname, cfg.Server.Port)

	return s.Serve()
}

func (s *server) Serve() error {
	return s.cli.Serve(s.ls)
}

func (s *server) Close() error {
	return s.ls.Close()
}
