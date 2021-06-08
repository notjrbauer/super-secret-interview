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

type StopCommand struct {
	cli proto.WorkerServiceClient
}

func NewStopCommand(cli proto.WorkerServiceClient) *StopCommand {
	return &StopCommand{
		cli: cli,
	}
}

func (c *StopCommand) Do(cmdAndArgs []string) error {
	if len(cmdAndArgs) == 0 {
		return errors.New("command length should be at least len(1)")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cmd := proto.StopRequest{
		JobID: cmdAndArgs[0],
	}
	_, err := c.cli.Stop(ctx, &cmd, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("stopping job %v\n", cmd.JobID))
	return nil
}
