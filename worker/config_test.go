package worker

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	var cfgBlob = `
  [global]
    name = "test_worker"
    log_dir = "/var/log/teleport/{{.Name}}"
  [server]
    hostname = "worker1.example.com"
    listen_addr = "127.0.0.1"
    listen_port = 6000
    ssl_cert = "/etc/teleport.d/worker1.cert"
    ssl_key = "/etc/teleport.d/worker1.key"
  `

	Convey("When giving invalid file", t, func() {
		cfg, err := LoadConfig("/path/to/invalid/file")
		So(err, ShouldNotBeNil)
		So(cfg, ShouldBeNil)
	})

	Convey("Everything should work on valid config file", t, func() {
		tmpfile, err := ioutil.TempFile("", "worker_test")
		So(err, ShouldEqual, nil)
		defer os.Remove(tmpfile.Name())

		tmpDir, err := ioutil.TempDir("", "worker_test")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tmpDir)

		curCfgBlob := cfgBlob

		err = ioutil.WriteFile(tmpfile.Name(), []byte(curCfgBlob), 0644)
		So(err, ShouldEqual, nil)
		defer tmpfile.Close()

		cfg, err := LoadConfig(tmpfile.Name())
		So(err, ShouldBeNil)
		So(cfg.Global.Name, ShouldEqual, "test_worker")

		So(cfg.Server.Hostname, ShouldEqual, "worker1.example.com")
	})

}
