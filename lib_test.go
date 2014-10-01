package pglib

import (
	"io"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testApi struct {
	UploadFileCalls int
}

func (ta *testApi) UploadFile(data FileData) (bool, error) {
	ta.UploadFileCalls++
	return true, nil
}

type BufferCloser struct {
	*io.PipeReader
	*io.PipeWriter
}

func (BufferCloser) Close() error {
	return nil
}

func TestRpcClient(t *testing.T) {
	p1r, p1w := io.Pipe()
	p2r, p2w := io.Pipe()
	ta := &testApi{}
	go ServeApi(ta, BufferCloser{p1r, p2w})

	a := ConnectApi(BufferCloser{p2r, p1w})
	Convey("test file upload call", t, func() {
		res, err := a.UploadFile(FileData{})
		So(err, ShouldEqual, nil)
		So(res, ShouldEqual, true)

		So(ta.UploadFileCalls, ShouldEqual, 1)
	})
}
