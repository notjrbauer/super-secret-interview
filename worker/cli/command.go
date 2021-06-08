package cli

import (
	"errors"
	"fmt"

	"github.com/notjrbauer/interview/super-secret-interview/worker"
)

type roundTripper interface {
	Do([]string) error
}

func Exec(cfg *worker.Config, cmdAndArgs []string) error {
	if len(cmdAndArgs) == 0 {
		return errors.New("command length should be at least len(1)")
	}
	cli, err := NewClient(cfg)
	if err != nil {
		return err
	}
	cmds := map[string]roundTripper{
		"stream": NewStreamCommand(cli),
		"start":  NewStartCommand(cli),
		"query":  NewQueryCommand(cli),
		"stop":   NewStopCommand(cli),
	}
	cmd, ok := cmds[cmdAndArgs[0]]
	if ok {
		return cmd.Do(cmdAndArgs[1:])
	}
	return fmt.Errorf("unknown command: %s", cmd)
}
