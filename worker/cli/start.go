package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/notjrbauer/interview/super-secret-interview/proto"
	"google.golang.org/grpc"
)

type StartCommand struct {
	cli proto.WorkerServiceClient
}

func NewStartCommand(cli proto.WorkerServiceClient) *StartCommand {
	return &StartCommand{
		cli: cli,
	}
}

func (c *StartCommand) Do(cmdAndArgs []string) error {
	if len(cmdAndArgs) == 0 {
		return errors.New("command length should be at least len(1)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	command := proto.StartRequest{ProcessName: cmdAndArgs[0]}
	command.Args = cmdAndArgs[1:]

	res, err := c.cli.Start(ctx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	os.Stdout.WriteString(fmt.Sprintf("starting job %v\n", res.JobID))
	return nil
}
