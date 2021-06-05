# RFC `APPROVED`: Job Worker Service Running(Linux Commands)

Editor: John Bauer
Status: **APPROVED**

### Requested reviewers: x, x, x, x.

### Approvals:

Team: **Cloud**

## Background

The assignment requires the implementation of a prototype job worker service
that provides an API to run arbitrary Linux processes. The assignment will
tackle Level 3.

## Technical Details && Criteria

The basic idea is to a have a simple client that communicates w/ a broker API, which enables the ability to query worker processes and the progress of their respected
jobs. The broker needs a persistent layer, so we will be using an `inmem sync(map[worker[jobid]])` to keep a reference of the worker and its progression w/ its job.

`jobID` will be generated with https://github.com/google/uuid

When a job initiates, three artifacts will be generated:

1. `/job/{uuid}` A job directory.
1. `/job/{uuid}/cmd.sh` file containing the raw command to execute.
1. `/job/{uuid}/log.txt` file containing the command's `stdout` and `stderr`.

We will be using `log.txt` as a source of record for streaming the output of
the job's executed command process.

### Worker Library

#### Methods

-   `Start` - Create a linux process and return the jobID && error.
-   `Stop` - Halts the job and in turn kills the job's running process. This
    will be issued via the `syscall.SIGTERM` signal to the command's process.
-   `Query` - Query for job to check status (Running, Exited).
    -   **TODO** explore if timestamps will be useful on the data object.
-   `Stream` - tails the output (stdout/stderr) of the executing process. (Make use of chans).

### API

-   GRPC API (start/stop/get status/stream output of a running process)
-   Responsible for providing authnz between client and server.

### Client

-   Should connect to worker service and schedule several jobs.
-   Should be able to query result of the job execution and stream its output.

## Design

Client -> Broker -> []PoolOfWorkers
![](https://github.com/donnemartin/system-design-primer/raw/master/images/h81n9iK.png)

## Security

### Transport Layer

We will be using Transport Layer Security (TLS). TLS stands for Transport Layer
Security and is the successor to SSL (Secure Sockets Layer). TLS provides secure
communication between web browsers and servers. The connection itself is secure
because symmetric cryptography is used to encrypt the data transmitted. The keys
are uniquely generated for each connection and are based on a shared secret
negotiated at the beginning of the session, also known as a TLS handshake.

We will be using `TLS v1.3`, which provides both privacy and data integrity in
securing communication between the client and the upstream server. `TLS v1.3` offers the following cipher suites:

-[x] TLS_AES_128_GCM_SHA256
-[] TLS_AES_256_GCM_SHA384
-[] TLS_CHACHA20_POLY_1305_SHA265
-[] TLS_AES_128_GCM_SHA256

#### Certs

-   X.509
-   Signature Algorithm: ecdsa-with-SHA256
-   Public Key Algorithm: id-ecPublicKey
-   ECDSA Public-Key (256bit)

### Authentication

We will be using mTLS for `Authentication`, which requires the generation of the
following artifacts:

#### CA

1. CA private key && self signed cert

#### Server

1. Private key && CSR
1. Signed cert from CA private key + CSR

#### Client

1. Private key && CSR
1. Signed cert from CA private key + CSR

### Authorization

A simple role based schema will be used for authorization:

1. Reader - query/stream
1. Writer - start/stop/query/stream

The client certificate contains an identifier on its `CN` (Common Name), which
will be mapped to its set of permissions, done through an imem map on the
backend.

```go
authorizedClients := map[string][]string{{
   "reading-client": []string{"read"},
   "writing-client": []string{"write"},
}}
...
methods := map[string][]string{{
  "query": []string{"read"},
  "stream": []string{"read"},
  "stop": []string{"write"},
  "start": []string{"write"},
}}

...
m := "stream"
clientIdentifier := certificate.CN
methodPermissions := method[m]
permissions, ok := authorizedClients[clientIdentifier]
if !ok {
 ...User Does not exist...
}
if ok := can(methodPermissions, permissions); !ok {
  ...User lacks authorization...
}

--Call Method Handler--
...
```

These roles will be mapped to the client's common name, or `CN`.

Using a [grpc interceptor](https://grpc.io/blog/grpc-web-interceptor/) allows us to check the original gRPC request's role context
before passing it along. We must support two types of actions: Request/Response (start,stop,query), and
Streaming data (stream). To support Request/Response, we will be using a Unary interceptor.
To support streaming, we will be using a [stream interceptor](https://grpc.io/blog/grpc-web-interceptor/#stream-interceptor-example) to check authz.

## Trade Offs

### Inmem vs Real Persistence

Restarting the broker service will wipe its queue, and the state of its old and new
jobs. This isn't ideal; however project scope makes an exception. Contrastingly,
disk persistence will allow us to amply query old jobs, and relay their valuable information to
the requesting client even after the main broker process has halted.

### Caching

No caching is required, (even inmem). We're trading off quick lookups, and fancy
LRU to handle large read load frequency. This is outside the scope of the project.

### Telemetry

## Definition of Success

-   High quality tests that cover happy and unhappy scenarios.
-   Project _should not be one giant pull request_
-   Program should compile and meet the technical details && criteria.
