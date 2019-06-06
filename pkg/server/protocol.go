package server

import (
	"bufio"
	"io"
)

const readChunkSize = 128

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

		_, err = writer.Write(command)
		if err != nil {
			panic("Failed to Write: " + err.Error())
		}
	}
}
