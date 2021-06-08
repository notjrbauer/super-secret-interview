package cli

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/notjrbauer/interview/super-secret-interview/proto"
	"github.com/notjrbauer/interview/super-secret-interview/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials(cfg *worker.Config) (credentials.TransportCredentials, error) {
	pemServerCA, err := ioutil.ReadFile(cfg.Server.CACert)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA %+v\n", pemServerCA)
	}

	clientCert, err := tls.LoadX509KeyPair(cfg.Client.SSLCert, cfg.Client.SSLKey)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
	}
	return credentials.NewTLS(tlsConfig), nil
}

func NewClient(cfg *worker.Config) (proto.WorkerServiceClient, error) {
	tlsCredentials, err := loadTLSCredentials(cfg)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.Server.Hostname, cfg.Server.Port),
		grpc.WithTransportCredentials(tlsCredentials),
	)
	if err != nil {
		return nil, err
	}
	return proto.NewWorkerServiceClient(conn), nil
}
