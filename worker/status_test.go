package worker

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	t.Run("Status json should work", func(t *testing.T) {
		b, err := json.Marshal(Success)
		assert.NoError(t, err)
		assert.Equal(t, []byte(`"success"`), b)

		var s JobStatus

		err = json.Unmarshal([]byte(`"failed"`), &s)
		assert.NoError(t, err)
		assert.Equal(t, Failed, s)
	})

}
