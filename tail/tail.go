package tail

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Stream accepts a file's path and returns channel which emits data from stdout
// and stderr.
func Stream(ctx context.Context, filePath string) (chan string, error) {
	ch := make(chan string)

	if ok := exists(filePath); !ok {
		return nil, errors.New("file descriptor does not exist")
	}
	cmd := exec.Command("tail", "-F", filePath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating stdout for Cmd", err)
	}

	// read from process' stdout, stderr and stdout are combined (&2>1).
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
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
		for {
			select {
			case <-ctx.Done():
				log.Printf("finished streaming for %s\n", filePath)
				return
			default:
				err := cmd.Wait()
				close(ch)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
				}
			}
		}
	}()

	return ch, nil
}

// Exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
