package server

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const httpGetReadChunkSize = 1024 * 1024

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
			fmt.Println("Failed to ReadBytes: " + err.Error())
			return
		}

		var shouldTerminate bool
		if command[len(command)-2] == '\r' {
			command[len(command)-2] = '\n'
			shouldTerminate = handleCommand(string(command[:len(command)-1]), writer)
		} else {
			shouldTerminate = handleCommand(string(command), writer)
		}

		if shouldTerminate {
			break
		}
	}
}

func reply(writer io.Writer, content []byte) bool {
	_, err := writer.Write(content)
	if err != nil {
		fmt.Println("Failed to Write: " + err.Error())
		return false
	}

	return true
}

func handleCommand(command string, writer io.Writer) bool {
	spaceIndex := strings.Index(command, " ")

	op := command[:len(command)-1]
	rest := ""
	if spaceIndex != -1 {
		op = command[:spaceIndex]
		rest = command[spaceIndex+1:]
	}

	var success bool
	switch op {
	case "echo":
		success = handlePing(rest, writer)
	case "get":
		success = handleGet(rest, writer)
	case "bye":
		return true
	default:
		success = reply(writer, []byte("unknown command"))
	}

	return !success
}

func handlePing(rest string, writer io.Writer) bool {
	return reply(writer, []byte(rest))
}

func handleGet(rest string, writer io.Writer) bool {
	resp, err := http.Get(rest[:len(rest)-1])
	if err != nil {
		fmt.Println("Failed to Get: " + err.Error())
		return false
	}

	buf := make([]byte, httpGetReadChunkSize)

	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Failed to Read: " + err.Error())
			return false
		}

		if !reply(writer, buf[:n]) {
			return false
		}

		if err == io.EOF {
			resp.Body.Close()
			return true
		}
	}
}
