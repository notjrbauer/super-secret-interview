package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/notjrbauer/interview/super-secret-interview/proto"
	"github.com/notjrbauer/interview/super-secret-interview/worker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func TestAuthedUser_AuthedMethod(t *testing.T) {
	var blob = `
  [global]
  name = "test_worker"

  [server]
  hostname = "localhost"
  listen_port = 6000
  ca_cert = "testdata/certs/root_ca.crt"
  ssl_cert = "testdata/certs/server.crt"
  ssl_key = "testdata/certs/server.key"

  [client]
  ca_cert = "testdata/certs/root_ca.crt"
  ssl_cert = "testdata/certs/client_write.crt"
  ssl_key = "testdata/certs/client_write.key"
  `
	config, err := createConfigFromBlob(t, blob)
	assert.NoError(t, err)

	serv := createTestServer(t, config)
	defer serv.Close()

	creds, err := loadClientCredentials(config)
	assert.NoError(t, err)

	target := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(creds))
	assert.NoError(t, err)
	defer func() {
		err := conn.Close()
		assert.NoError(t, err)
	}()

	client := proto.NewWorkerServiceClient(conn)
	res, err := client.Start(context.Background(), &proto.StartRequest{ProcessName: "ls"})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.JobID)
}

func TestAuthedUser_UnauthzMethod(t *testing.T) {
	var blob = `
  [global]
  name = "test_worker"

  [server]
  hostname = "localhost"
  listen_port = 6000
  ca_cert = "testdata/certs/root_ca.crt"
  ssl_cert = "testdata/certs/server.crt"
  ssl_key = "testdata/certs/server.key"

  [client]
  ca_cert = "testdata/certs/root_ca.crt"
  ssl_cert = "testdata/certs/client_unauthorized.crt"
  ssl_key = "testdata/certs/client_unauthorized.key"
  `
	config, err := createConfigFromBlob(t, blob)
	assert.NoError(t, err)

	serv := createTestServer(t, config)
	defer serv.Close()

	creds, err := loadClientCredentials(config)
	assert.NoError(t, err)

	target := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(creds))
	assert.NoError(t, err)

	defer func() {
		err := conn.Close()
		assert.NoError(t, err)
	}()

	client := proto.NewWorkerServiceClient(conn)

	res, err := client.Start(context.Background(), &proto.StartRequest{ProcessName: "ls"})
	assert.Error(t, err)
	assert.Nil(t, res)

	code, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "unauthorized", code.Message())
	assert.Nil(t, res)
}

func TestUntrustedUser(t *testing.T) {
	var blob = `
  [global]
  name = "test_worker"

  [server]
  hostname = "localhost"
  listen_port = 6000
  ca_cert = "testdata/certs/root_ca.crt"
  ssl_cert = "testdata/certs/server.crt"
  ssl_key = "testdata/certs/server.key"

  [client]
  ca_cert = "testdata/certs/evil_ca.crt"
  ssl_cert = "testdata/certs/client_untrusted.crt"
  ssl_key = "testdata/certs/client_untrusted.key"
  `
	config, err := createConfigFromBlob(t, blob)
	assert.NoError(t, err)

	serv := createTestServer(t, config)
	defer serv.Close()

	creds, err := loadClientCredentials(config)
	assert.NoError(t, err)

	target := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(creds))
	assert.NoError(t, err)

	defer func() {
		err := conn.Close()
		assert.NoError(t, err)
	}()

	client := proto.NewWorkerServiceClient(conn)

	res, err := client.Start(context.Background(), &proto.StartRequest{ProcessName: "ls"})
	assert.Nil(t, res)
	assert.Error(t, err)

	code, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unavailable, code.Code())
}

func createTestServer(t *testing.T, cfg *worker.Config) *server {
	creds, err := loadTLSCredentials(cfg)
	assert.NoError(t, err)

	serv, err := newServer(cfg, creds)
	assert.NoError(t, err)

	go serv.Serve()
	return serv
}

func createConfigFromBlob(t *testing.T, blob string) (*worker.Config, error) {
	tmpfile, err := ioutil.TempFile("", "worker_test")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())

	curCfgBlob := blob

	err = ioutil.WriteFile(tmpfile.Name(), []byte(curCfgBlob), 0644)
	if err != nil {
		return nil, err
	}
	defer tmpfile.Close()

	cfg, err := worker.LoadConfig(tmpfile.Name())
	cfg.Global.LogDir = t.TempDir()

	return cfg, err
}

func loadClientCredentials(cfg *worker.Config) (credentials.TransportCredentials, error) {
	pemClientCA, err := ioutil.ReadFile(cfg.Client.CACert)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, errors.New("failed to add CA's certificate")
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
