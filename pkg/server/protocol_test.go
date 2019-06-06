package server

import (
	"bytes"
	"testing"
)

func TestProtocolEcho(t *testing.T) {
	input := bytes.NewBufferString("hello, world!")
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte("hello, world!")) {
		t.Fatalf("Unexpected output: %s", output.Bytes())
	}
}
