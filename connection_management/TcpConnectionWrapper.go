package connection_management

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const MessageLengthSize = 10

func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

type TcpConnectionWrapper struct {
	net.Conn
}

func (conn *TcpConnectionWrapper) ReadMessage() ([]byte, error) {
	buffer := make([]byte, MessageLengthSize)

	if _, err := io.ReadAtLeast(io.LimitReader(conn, MessageLengthSize), buffer, MessageLengthSize); err != nil {
		return nil, err
	}

	msgLength, err := strconv.Atoi(string(buffer[:MessageLengthSize]))
	if err != nil {
		return nil, err
	}

	outputBuffer := make([]byte, msgLength)
	bytesRead := 0

	for bytesRead < msgLength {
		n, err := conn.Read(outputBuffer[bytesRead:])
		if err != nil {
			return nil, err
		}
		bytesRead += n
	}
	return outputBuffer, nil
}

func (conn *TcpConnectionWrapper) ReadMessageToFile(file *os.File) error {
	buffer := make([]byte, MessageLengthSize)

	if _, err := io.ReadAtLeast(io.LimitReader(conn, MessageLengthSize), buffer, MessageLengthSize); err != nil {
		return err
	}

	msgLength, err := strconv.Atoi(string(buffer[:MessageLengthSize]))
	if err != nil {
		return err
	}

	buffer = make([]byte, 2048)
	bytesRead := 0

	for bytesRead < msgLength {
		n, err := conn.Read(buffer)
		if err != nil {
			return err
		}
		bytesRead += n

		_, err = file.Write(buffer[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func (conn *TcpConnectionWrapper) WriteMessage(msg []byte) error {
	_, err := conn.Write([]byte(leftPad2Len(fmt.Sprint(len(msg)), "0", 10)))
	if err != nil {
		return err
	}
	_, err = conn.Write(msg)
	return err
}

func (conn *TcpConnectionWrapper) WriteFileMessage(file *os.File) error {
	buffer := make([]byte, 2048)

	stats, _ := file.Stat()
	size := leftPad2Len(fmt.Sprint(stats.Size()), "0", 10)
	if _, err := conn.Write([]byte(size)); err != nil {
		return err
	}

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if _, err = conn.Write(buffer[:n]); err != nil {
			return err
		}
	}

	return nil
}

func (conn *TcpConnectionWrapper) WriteStatusCode(statusCode int) error {
	return conn.WriteMessage([]byte(fmt.Sprint(statusCode)))
}
