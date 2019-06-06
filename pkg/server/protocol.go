package server

import "io"

const readChunkSize = 128

// HandleProtocol handles the protocol part of a simple proxy connection. It
// reads commands from reader, and writes output to writer.
func HandleProtocol(reader io.Reader, writer io.Writer) {
	readBuf :=
		make([]byte, readChunkSize)

	for {
		n, err := reader.Read(readBuf)
		if err != nil && err != io.EOF {
			panic("Failed to read: " + err.Error())
		}

		if n == 0 {
			return
		}

		_, err = writer.Write(readBuf[:n])
		if err != nil {
			panic("Failed to write: " + err.Error())
		}
	}
}
