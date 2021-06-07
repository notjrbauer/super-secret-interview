package worker

import (
	"encoding/json"
	"errors"
	"fmt"
)

type JobStatus uint8

const (
	None JobStatus = iota
	Failed
	Success
	Running
	Stopped
)

func (s JobStatus) String() string {
	switch s {
	case None:
		return "none"
	case Failed:
		return "failed"
	case Success:
		return "success"
	case Running:
		return "running"
	case Stopped:
		return "stopped"
	default:
		return "unknown"
	}
}

func (s JobStatus) MarshalJSON() ([]byte, error) {
	strStatus := s.String()
	if strStatus == "" {
		return []byte{}, errors.New("Invalid status value")
	}

	return json.Marshal(strStatus)
}

func (s *JobStatus) UnmarshalJSON(v []byte) error {
	sv := string(v)
	switch sv {
	case `"none"`:
		*s = None
	case `"failed"`:
		*s = Failed
	case `"success"`:
		*s = Success
	case `"running"`:
		*s = Running
	default:
		return fmt.Errorf("Invalid status value: %s", string(v))
	}
	return nil
}
