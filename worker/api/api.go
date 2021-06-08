package api

import (
	"context"

	"github.com/notjrbauer/interview/super-secret-interview/proto"
	"github.com/notjrbauer/interview/super-secret-interview/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type workerServer struct {
	proto.UnimplementedWorkerServiceServer
	Service worker.Service
}

// Start starts a process.
func (s *workerServer) Start(ctx context.Context, r *proto.StartRequest) (*proto.StartResponse, error) {

	jobID, err := s.Service.Start(append([]string{r.ProcessName}, r.Args...))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := proto.StartResponse{
		JobID: jobID,
	}
	return &out, nil
}

// Stop stops a process.
func (s *workerServer) Stop(ctx context.Context, r *proto.StopRequest) (*proto.StopResponse, error) {

	err := s.Service.Stop(r.JobID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.StopResponse{}, nil
}

// Query performs a lookup on currently running / past run processes by the manager.
func (s *workerServer) Query(ctx context.Context, r *proto.QueryRequest) (*proto.QueryResponse, error) {

	job, err := s.Service.Query(r.JobID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := proto.QueryResponse{
		ProcessID: int64(job.PID),
		ExitCode:  int64(job.ExitCode),
		Status:    proto.QueryResponse_Status(job.Status),
	}
	return &out, nil
}

// Stream streams a resource specified in the request.
func (s *workerServer) Stream(r *proto.StreamRequest, stream proto.WorkerService_StreamServer) error {

	eventCh, err := s.Service.Stream(stream.Context(), r.JobID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case eventData, ok := <-eventCh:
			if !ok {
				return nil
			}

			out := &proto.StreamResponse{Chunk: eventData}

			if err := stream.SendMsg(out); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}
}
