package main_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/berquerant/firehose-proto/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestRun(t *testing.T) {
	dir := newTmpDir("grpcx")
	defer dir.close()
	bin := dir.path("server")
	assert.Nil(t, newCommand("go", "build", "-o", bin).Run(), "build binary")

	const port = "10001"
	os.Setenv("PORT", port)

	runCommand := newCommand(bin)
	assert.Nil(t, runCommand.Start(), "start server")

	time.Sleep(time.Second) // wait for server to start up

	func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("127.0.0.1:%s", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		assert.Nil(t, err, "connect to server")
		defer conn.Close()

		client := empty.NewEmptyServiceClient(conn)
		_, err = client.Ping(context.TODO(), new(emptypb.Empty))
		assert.Nil(t, err, "ping")
	}()

	assert.Nil(t, runCommand.Process.Signal(os.Interrupt), "sigint")
	assert.Nil(t, runCommand.Wait(), "wait for server to stop")
}

func newCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

type tmpDir struct {
	dir string
}

func newTmpDir(name string) *tmpDir {
	dir, err := os.MkdirTemp("", name)
	if err != nil {
		panic(err)
	}
	return &tmpDir{
		dir: dir,
	}
}

func (t *tmpDir) path(name string) string {
	return filepath.Join(t.dir, name)
}

func (t *tmpDir) close() {
	os.RemoveAll(t.dir)
}
