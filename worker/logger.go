package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Logger mounts a directory w/ and generates a log file.
type Logger struct {
	LogDir string
}

// NewLogger returns a new Logger instance.
func NewLogger(cfg *Config) *Logger {
	return &Logger{
		LogDir: cfg.Global.LogDir,
	}
}

// Path returns a absolute file path given a name.
func (l *Logger) Path(name string) string {
	return filepath.Join(l.LogDir, fmt.Sprintf("%s", name))
}

// Create creates the directory of the given path w/ the logger's base path. It then creates the file specified by the fileName.
func (l *Logger) Create(path string, fileName string) (*os.File, error) {
	path = l.Path(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Making dir %s\n", path)
		if err = os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("Error making dir %s: %w", path, err)
		}
	}
	return os.Create(fmt.Sprintf("%s/%s", path, fileName))
}

// Remove removes the specified file.
func (l *Logger) Remove(name string) error {
	return os.Remove(l.Path(name))
}

// Stream accepts a file's name and returns channel which emits data from the tail'd file.
func (l *Logger) Stream(ctx context.Context, filePath string, fileName string) (chan string, error) {
	ch := make(chan string)

	filePath = l.Path(filePath)
	filePath = fmt.Sprintf("%s/%s", filePath, fileName)
	if ok := exists(filePath); !ok {
		return nil, errors.New("file descriptor does not exist")
	}

	cmd := exec.Command("tail", "-F", filePath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating stdout for Cmd", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		// todo: this should be configurable w/ a limit
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading from stdout for Cmd", err)
				return
			}
			ch <- string(buf[0:n])
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			log.Printf("finished streaming for %s\n", filePath)
			return
		default:
			err := cmd.Wait()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
			}
			close(ch)
		}
	}()

	return ch, nil
}

// exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
