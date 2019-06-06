package server

import (
	"bufio"
	"io"
	"net/http"
	"strings"
)

const httpGetReadChunkSize = 1024

// HandleProtocol handles the protocol part of a simple proxy connection. It
// reads commands from reader, and writes output to writer.
func HandleProtocol(reader io.Reader, writer io.Writer) {
	bufReader := bufio.NewReader(reader)

	for {
		command, err := bufReader.ReadBytes('\n')
		if err == io.EOF || command[len(command)-1] != '\n' {
			return
		}
		if err != nil {
			panic("Failed to ReadBytes: " + err.Error())
		}

		handleCommand(string(command), writer)
	}
}

func reply(writer io.Writer, content []byte) {
	_, err := writer.Write(content)
	if err != nil {
		panic("Failed to Write: " + err.Error())
	}
}

func handleCommand(command string, writer io.Writer) {
	spaceIndex := strings.Index(command, " ")
	if spaceIndex == -1 {
		reply(writer, []byte("ill-formatted command"))
		return
	}

	op := command[:spaceIndex]
	rest := command[spaceIndex+1:]

	switch op {
	case "echo":
		handlePing(rest, writer)
	case "get":
		handleGet(rest, writer)
	default:
		reply(writer, []byte("unknown command"))
	}
}

func handlePing(rest string, writer io.Writer) {
	reply(writer, []byte(rest))
}

func handleGet(rest string, writer io.Writer) {
	resp, err := http.Get(rest[:len(rest)-1])
	if err != nil {
		panic("Failed to Get: " + err.Error())
	}

	buf := make([]byte, httpGetReadChunkSize)

	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			panic("Failed to Read: " + err.Error())
		}

		reply(writer, buf[:n])
		if err == io.EOF {
			resp.Body.Close()
			return
		}
	}
}
