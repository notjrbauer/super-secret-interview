package cli

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/notjrbauer/interview/super-secret-interview/proto"
	"google.golang.org/grpc"
)

type StreamCommand struct {
	cli proto.WorkerServiceClient
}

func NewStreamCommand(cli proto.WorkerServiceClient) *StreamCommand {
	return &StreamCommand{
		cli: cli,
	}
}

func (c *StreamCommand) Do(cmdAndArgs []string) error {
	if len(cmdAndArgs) == 0 {
		return errors.New("command length should be at least len(1)")
	}

	cmd := proto.StreamRequest{}
	cmd.JobID = cmdAndArgs[0]

	ctx, cancel := context.WithCancel(context.Background())
	res, err := c.cli.Stream(ctx, &cmd, grpc.WaitForReady(true))
	if err != nil {
		cancel()
		return err
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	defer func() {
		cancel()
		signal.Stop(done)
	}()

	go func() {
		for {
			out, err := res.Recv()
			if err != nil {
				return
			}
			os.Stdout.WriteString(out.Chunk)
		}
	}()

	<-done
	return nil
}
