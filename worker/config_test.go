package worker

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

	t.Run("When giving an invalid file", func(t *testing.T) {
		cfg, err := LoadConfig("/path/to/invalid/file")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("Everything should work on valid config file", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "worker_test")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		tmpDir, err := ioutil.TempDir("", "worker_test")
		assert.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		curCfgBlob := cfgBlob

		err = ioutil.WriteFile(tmpfile.Name(), []byte(curCfgBlob), 0644)
		assert.NoError(t, err)
		defer tmpfile.Close()

		cfg, err := LoadConfig(tmpfile.Name())
		assert.NoError(t, err)
		assert.Equal(t, "test_worker", cfg.Global.Name)
		assert.Equal(t, "worker1.example.com", cfg.Server.Hostname)
	})
}
