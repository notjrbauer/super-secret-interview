package worker

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStatus(t *testing.T) {
  Convey("Status json should work", t, func() {

		b, err := json.Marshal(Success)
		So(err, ShouldBeNil)
		So(b, ShouldResemble, []byte(`"success"`)) // deep equal should be used

		var s JobStatus

		err = json.Unmarshal([]byte(`"failed"`), &s)
		So(err, ShouldBeNil)
		So(s, ShouldEqual, Failed)
  })
}
