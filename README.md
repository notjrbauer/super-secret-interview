# Teleport Interview

## Build

```bash
# Build API and CLI
make all
```

## Test

```sh
make test
```

## Run API

```sh
./bin/api --config worker.conf
```

## Run CLI

```sh
./bin/cli --config worker.conf
```

## Configuration

Configuration is defined below. Certificates are available with different permissions and are located in `/certs`.

```sh
◗ tree certs
certs
├── client_read.crt
├── client_read.csr
├── client_read.key
├── client_unauthorized.crt
├── client_unauthorized.csr
├── client_unauthorized.key
├── client_unauthorized_permissions.crt
├── client_unauthorized_permissions.csr
├── client_unauthorized_permissions.key
├── client_untrusted.crt
├── client_untrusted.csr
├── client_untrusted.key
├── client_write.crt
├── client_write.csr
├── client_write.key
├── evil_ca.crt
├── evil_ca.key
├── root_ca.crt
├── root_ca.key
├── server.crt
├── server.csr
└── server.key
0 directories, 22 files
```

```toml
# vim: ft=toml

[global]
name = "test_worker"
log_dir = "/tmp/teleport"

[server]
hostname = "localhost"
listen_addr = "127.0.0.1"
listen_port = 6000

ca_cert = "certs/root_ca.crt"
ssl_cert = "certs/server.crt"
ssl_key = "certs/server.key"

[client]
ca_cert = "certs/root_ca.crt"
ssl_cert = "certs/client_write.crt"
ssl_key = "certs/client_write.key"
```

### Usage

```sh
# Start a process!
❮ ./bin/cli start "bash" "-c" "until((0)); do date; done"
starting job c7c75447-1ea7-40a1-add4-f89ffb8a921d

# Query a response!
ϟ ./bin/cli query c7c75447-1ea7-40a1-add4-f89ffb8a921d
Query Response: processID:83614 status:RUNNING+

# Stream a process log file!
ϟ  ./bin/cli stream c7c75447-1ea7-40a1-add4-f89ffb8a921d

Tue Jun  8 01:31:18 EDT 2021
Tue Jun  8 01:31:18 EDT 2021
Tue Jun  8 01:31:18 EDT 2021
Tue Jun  8 01:31:18 EDT 2021
Tue Jun  8 01:31:18 EDT 2021
Tue Jun  8 01:31:18 EDT 2021
Tue Jun  8 01:31:18 EDT 2021

# Stop a process!
ϟ  ./bin/cli stop c7c75447-1ea7-40a1-add4-f89ffb8a921d

stopping job c7c75447-1ea7-40a1-add4-f89ffb8a921d

# Query historical records!
ϟ ./bin/cli query c7c75447-1ea7-40a1-add4-f89ffb8a921d

Query Response: processID:83614 exitCode:-1 status:STOPPED+

```
