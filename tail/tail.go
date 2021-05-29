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
func Stream(context context.Context, filePath string) (chan string, error) {
	ch := make(chan string)

	if ok := Exists(filePath); !ok {
		return nil, errors.New("file descriptor does not exist")
	}
	go func() {
		cmd := exec.Command("tail", "-F", filePath)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating stdout for Cmd", err)
		}

		cmd.Stderr = cmd.Stdout
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		// read from process' stdin, stderr and stdin are combined (&2>1).
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			ch <- string(buf[0:n])
		}

		close(ch)

		if err := cmd.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		}
	}()

	log.Printf("finished streaming for %s\n", filePath)

	return ch, nil
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
