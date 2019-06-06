package server

import (
	"bytes"
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

func TestIllFormattedCommand(t *testing.T) {
	input := bytes.NewBufferString("test\n")
	output := bytes.Buffer{}

	HandleProtocol(input, &output)

	if !bytes.Equal(output.Bytes(), []byte("ill-formatted command")) {
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
