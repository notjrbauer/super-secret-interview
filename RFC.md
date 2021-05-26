# RFC 1 `UNAPPROVED`: Job Worker Service Running(Linux Commands)

Editor: John Bauer
Status: **UNAPPROVED**

### Requested reviewers: x, x, x, x.

### Approvals:

Team: **Cloud**

## Background

The assignment requires the implementation of a prototype job worker service
that provides an API to run arbitrary Linux processes. The assignment will
tackle Level 3.

## Technical details && Criteria

The basic idea is to a have a simple client that communicates w/ a broker API, which enables the ability to query worker processes and the progress of their respected
jobs. The broker needs a persistent layer, so we will be using an `inmem sync(map[worker[jobid]])` to keep a reference of the worker and its progression w/ its job.

### Worker Library

-   Start / Stop / Query Status of a job.
-   Stream output of a running job.

### API

-   GRPC API (start/stop/get status/stream output of a running process)

```go

type Command struct {
 UUID string
 Command string
 Arguments []string
 User string
 Group string
 Directory string
}

type Message struct {
  line uint32
  message string
}

type JobRequest struct {
  Name string

  Command
}

type JobResponse struct {
  Name string
  Duration time.Duration

  Command
  Message
}
```

### Client

-   Should connect to worker service and schedule several jobs.
-   Should be able to query result of the job execution and stream its output.

## Design

Client -> Broker -> []PoolOfWorkers
![](https://github.com/donnemartin/system-design-primer/raw/master/images/h81n9iK.png)

## Trade Offs

### Syncing State - Channels vs Mutex

Channels will allow an easier primative for allowing the broker routine to communicate/keep track of our worker routines.

### Inmem vs Real Persistence

Restarting the broker service will wipe its queue, and the state of its old and new
jobs. This isn't ideal; however project scope makes an exception. Contrastingly,
disk persistence will allow us to amply query old jobs, and relay their valuable information to
the requesting client even after the main broker process has halted.

### Caching

No caching is required, (even inmem). We're trading off quick lookups, and fancy
LRU to handle large read load frequency. This is outside the scope of the project.

### Telemetry

## Security

-   Authentication: Use mTLS and verify client cert. Setup strong cipher suites for TLS and
    good crpyto certs.
-   Authorization: Use simple authz scheme.

## Definition of Success

-   High quality tests that cover happy and unhappy scenarios.
-   Project _should not be one giant pull request_
-   Program should compile and meet the technical details && criteria.

```

```
