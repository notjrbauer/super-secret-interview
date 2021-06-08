package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCan(t *testing.T) {
	assert.True(t, Can("/WorkerService/Start", "client_write"))
	type test struct {
		name   string
		method string
		user   string
		want   bool
	}

	tests := []test{
		{
			name:   "User Matches Method and Permission",
			method: "/WorkerService/Start",
			user:   "client_write",
			want:   true,
		},
		{
			name:   "User Matches Method but not Permission",
			method: "/WorkerService/Start",
			user:   "client_read",
			want:   false,
		},
		{
			name:   "User Matches Method but invalid User Role",
			method: "/WorkerService/Start",
			user:   "client_unauthorized_permissions",
			want:   false,
		},
		{
			name:   "User Matches but has invalid method",
			method: "/WorkerService/InvalidMethod",
			user:   "client_unauthorized_permissions",
			want:   false,
		},
	}

	for _, tc := range tests {
		got := Can(tc.method, tc.user)
		assert.Equal(t, tc.want, got, tc.name)
	}
}
