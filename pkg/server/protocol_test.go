package server

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"testing"
)

func TestProtocolEcho(t *testing.T) {
	input := bytes.NewBufferString("echo hello, world!\n")
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte("hello, world!\n")) {
		t.Fatalf("Unexpected output: %s", output.Bytes())
	}
}

func TestIncompleteLine(t *testing.T) {
	const testString = "hello, world!"
	input := bytes.NewBufferString(testString)
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte("")) {
		t.Fatalf("Unexpected output: %s", output.Bytes())
	}
}

func TestUnknownCommand(t *testing.T) {
	input := bytes.NewBufferString("test \n")
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte("unknown command")) {
		t.Fatalf("Unexpected output: %s", output.Bytes())
	}
}

type FixedReplyHTTPHandler struct {
	reply string
}

func (handler FixedReplyHTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	reply(writer, []byte(handler.reply))
}

func TestGetOK(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic("Failed to Listen: " + err.Error())
	}

	port := listener.Addr().(*net.TCPAddr).Port

	go func() {
		err := http.Serve(listener, FixedReplyHTTPHandler{"test-content"})
		if err != nil {
			panic("Failed to Serve: " + err.Error())
		}
	}()

	input := bytes.NewBufferString(fmt.Sprintf("get http://localhost:%d\n", port))
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte("test-content")) {
		t.Fatalf("Unexpected output: %s", output.Bytes())
	}
}

func TestCrLfLineEnding(t *testing.T) {
	input := bytes.NewBufferString("echo hello, world!\r\n")
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte("hello, world!\n")) {
		t.Fatalf("Unexpected output: %s", output.Bytes())
	}
}

func TestBye(t *testing.T) {
	input := bytes.NewBufferString("echo a\nbye\n echo b\n")
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte("a\n")) {
		t.Fatalf("Unexpected output: %s", output.Bytes())
	}
}
