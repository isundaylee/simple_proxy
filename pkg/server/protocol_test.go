package server

import (
	"bytes"
	"testing"
)

func TestProtocolEcho(t *testing.T) {
	const testString = "hello, world!\n"
	input := bytes.NewBufferString(testString)
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte(testString)) {
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
