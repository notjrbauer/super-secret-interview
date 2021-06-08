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

type QueryCommand struct {
	cli proto.WorkerServiceClient
}

func NewQueryCommand(cli proto.WorkerServiceClient) *QueryCommand {
	return &QueryCommand{
		cli: cli,
	}
}

func (c *QueryCommand) Do(cmdAndArgs []string) error {
	if len(cmdAndArgs) == 0 {
		return errors.New("command length should be at least len(1)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	cmd := proto.QueryRequest{}
	cmd.JobID = cmdAndArgs[0]

	res, err := c.cli.Query(ctx, &cmd, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	os.Stdout.WriteString(fmt.Sprintf("Query Response: %v+\n", res))
	return nil
}
