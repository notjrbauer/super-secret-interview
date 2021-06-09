package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

const terminationGraceInterval = (time.Second * 2)

var errProcessNotStarted = errors.New("Process Not Started")

// Service interacts with linux processes.
type Service interface {

	// Start creates a linux process.
	Start(cmd []string) (jobID string, err error)
	// Stop kills a linux process.
	Stop(jobID string) (err error)
	// Query returns the status of a cmd process.
	Query(jobID string) (status CmdStatus, err error)
	// Streams the process output.
	Stream(ctx context.Context, jobID string) (emitter chan string, err error)
}

// CmdStatus represents the status of a cmd process.
type CmdStatus struct {
	// PID is the process id.
	PID int
	// Status is an enum mapping to the states above.
	Status JobStatus
	// ExitCode represents the exit code of the process, or -1 if it's still running.
	ExitCode int
}

type empty struct{}

// Job is a linux process launched by the worker.
type Job struct {
	//  ID is the unique identifying ID of the job.
	ID string
	// Cmd is the command executor.
	Cmd *exec.Cmd
	// Status is the process status.
	Status *CmdStatus
	// WorkingDir is the working directory of the job
	WorkingDir string
	// logFileFd is the logfile file descriptor
	logFileFd *os.File
	// finished
	finished chan empty
}

// NewJob returns a new cmd job to run
func NewJob(cmdAndArgs []string) (*Job, error) {
	var cmd *exec.Cmd
	jobID := uuid.NewString()

	if len(cmdAndArgs) == 0 {
		return nil, errors.New("command length should be at least len(1)")
	}
	cmd = exec.Command(cmdAndArgs[0], cmdAndArgs[1:]...)

	return &Job{Cmd: cmd, ID: jobID}, nil
}

// Terminate issues a SIGINT, and then a SIGTERM to kill the job
func (j *Job) Terminate() error {
	if j.Cmd == nil || j.Cmd.Process == nil {
		return errProcessNotStarted
	}

	if err := j.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("error killing process %d : %w", j.Cmd.Process.Pid, err)
	}
	select {
	case <-time.After(terminationGraceInterval):
		j.Cmd.Process.Signal(syscall.SIGKILL)
		log.Println("SIGTERM failed to kill the job in 2s. SIGKILL sent")
	case <-j.finished:

	}
	return nil
}

// SetLogFile binds a file descriptor to a process' stdout, stderr.
func (j *Job) SetLogFile(logFile *os.File) {
	j.Cmd.Stdout = logFile
	j.Cmd.Stderr = logFile
	j.logFileFd = logFile
}

// NewService returns a new worker service.
func NewService(cfg *Config) *workerService {
	return &workerService{
		cfg:  cfg,
		log:  NewLogger(cfg),
		jobs: make(map[string]*Job),
	}
}

type workerService struct {
	mu sync.RWMutex

	cfg  *Config
	log  *Logger
	jobs map[string]*Job
}

// Start creates a process with the command and arguments supplied.
func (w *workerService) Start(cmdAndArgs []string) (string, error) {
	job, err := NewJob(cmdAndArgs)
	if err != nil {
		return "", fmt.Errorf("error starting job: %w", err)
	}

	logFileFd, err := w.log.Create(job.ID, "log_txt")
	if err != nil {
		return "", err
	}

	job.SetLogFile(logFileFd)

	job.finished = make(chan empty)

	if err := job.Cmd.Start(); err != nil {
		w.log.Remove(job.ID)
		return job.ID, err
	}

	job.Status = &CmdStatus{PID: job.Cmd.Process.Pid, Status: Running}
	w.mu.Lock()
	w.jobs[job.ID] = job
	w.mu.Unlock()

	go func() {
		defer func() { job.finished <- empty{} }()
		if err := job.Cmd.Wait(); err != nil {
			log.Printf("command execution failed: %v", err)
		}

		updatedStatus := CmdStatus{
			Status:   Stopped,
			PID:      job.Cmd.ProcessState.Pid(),
			ExitCode: job.Cmd.ProcessState.ExitCode(),
		}

		code := updatedStatus.ExitCode
		isExited := job.Cmd.ProcessState.Exited()

		if isExited {
			if code > 0 {
				updatedStatus.Status = Failed
			}
			if code == 0 {
				updatedStatus.Status = Success
			}
		}

		w.mu.Lock()
		job.Status = &updatedStatus
		w.mu.Unlock()
	}()

	return job.ID, nil
}

// Stop sends a SIGTERM to the job process. After a 2 second timeout,
// the process is issued the SIGKILL signal.
func (w *workerService) Stop(jobID string) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	job, ok := w.jobs[jobID]
	if !ok {
		return fmt.Errorf("job not found: %s", jobID)
	}

	if err := job.Terminate(); err != nil {
		return fmt.Errorf("error terminating jobID:%s, %w", jobID, err)
	}

	return nil
}

// Query returns the CmdStatus of a job.
func (w *workerService) Query(jobID string) (CmdStatus, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	job, ok := w.jobs[jobID]
	if !ok {
		return CmdStatus{}, fmt.Errorf("job not found: %s", jobID)
	}

	return *job.Status, nil
}

// Stream reads from the log file, like 'tail -f' through
// a channel. If the context is canceled the channel will
// be closed and the tailing will be stopped.
func (w *workerService) Stream(ctx context.Context, jobID string) (chan string, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	job, ok := w.jobs[jobID]
	if !ok {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}
	// TODO: Make this a const for the log file.
	return w.log.Stream(ctx, job.ID, "log_txt")
}
