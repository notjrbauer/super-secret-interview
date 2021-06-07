package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkerService_ShortProcess(t *testing.T) {
	c := &config{}
	c.Global.LogDir = t.TempDir()
	jobID := ""

	ws := NewService(c)
	t.Run("Basic Command", func(t *testing.T) {
		var err error
		jobID, err = ws.Start([]string{"echo", "determinism"})
		assert.NoError(t, err)
		assert.NotEmpty(t, jobID)
	})
	t.Run("Query Proccess_Success", func(t *testing.T) {
		s, err := ws.Query(jobID)
		assert.NoError(t, err)
		assert.Equal(t, Running, s.Status)
	})
	t.Run("Stream Proccess_Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		logEvents, err := ws.Stream(ctx, jobID)
		assert.NoError(t, err)
		assert.Equal(t, "determinism\n", <-logEvents)
		cancel()
	})
	t.Run("Stop Already Successful Proccess_Success", func(t *testing.T) {
		s, err := ws.Query(jobID)
		assert.NoError(t, err)
		assert.Equal(t, Success, s.Status)

		err = ws.Stop(jobID)
		assert.Error(t, err)
	})
}

func TestNewWorkerService_LongProcess(t *testing.T) {
	c := &config{}
	c.Global.LogDir = t.TempDir()
	jobID := ""

	ws := NewService(c)
	t.Run("Basic Command", func(t *testing.T) {
		var err error
		jobID, err = ws.Start([]string{"bash", "-c", "until ((0)); do date; sleep 2; done"})
		assert.NoError(t, err)
		assert.NotEmpty(t, jobID)
	})

	t.Run("Query Proccess_Success", func(t *testing.T) {
		s, err := ws.Query(jobID)
		assert.NoError(t, err)
		assert.Equal(t, Running, s.Status)
	})

	t.Run("Stream Proccess_Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		logEvents, err := ws.Stream(ctx, jobID)
		assert.NoError(t, err)
		assert.NotNil(t, <-logEvents)
		assert.NotNil(t, <-logEvents)
		cancel()
	})

	t.Run("Stop Proccess_Success", func(t *testing.T) {
		s, err := ws.Query(jobID)
		assert.NoError(t, err)
		assert.Equal(t, Running, s.Status)

		err = ws.Stop(jobID)
		assert.NoError(t, err)

		time.Sleep(time.Second * 1)

		s, err = ws.Query(jobID)
		assert.NoError(t, err)
		assert.Equal(t, Stopped, s.Status)
	})
}
