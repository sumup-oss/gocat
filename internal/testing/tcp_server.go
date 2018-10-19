// Copyright 2018 SumUp Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testing

import (
	"io"
	"net"

	"github.com/stretchr/testify/require"
)

type TCPServer struct {
	address         string
	t               TestingT
	msgsToBroadcast chan []byte
	bufferSize      int
}

func NewTCPServer(t TestingT, bufferSize int, address string) *TCPServer {
	return &TCPServer{
		address:         address,
		t:               t,
		bufferSize:      bufferSize,
		msgsToBroadcast: make(chan []byte, 1000),
	}
}

func (ts *TCPServer) Serve(started chan<- *ListenResult) {
	ln, err := net.Listen("tcp", ts.address)
	if err != nil {
		started <- &ListenResult{
			Err: err,
		}

		return
	}
	defer ln.Close()

	addr := ln.Addr().String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		started <- &ListenResult{
			Address: addr,
			Err:     err,
		}
	}

	started <- &ListenResult{
		Address: addr,
		Port:    port,
		Host:    host,
	}

	for {
		c, err := ln.Accept()
		require.Nil(ts.t, err, "Failed to accept connection")

		go ts.handleConnection(c)
	}
}

func (ts *TCPServer) handleConnection(c net.Conn) {
	defer c.Close()

	buf := make([]byte, ts.bufferSize)

	for {
		n, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}

			ts.t.Logf("Copy Err: %s", err.Error())
			return
		}

		if n < 1 {
			break
		}

		msg := buf[:n]
		writtenBytes, err := c.Write(msg)
		require.Nil(ts.t, err, "Failed to write data")

		expectedWrittenBytes := len(msg)
		if expectedWrittenBytes != writtenBytes {
			ts.t.Fatalf(
				"Incomplete write back, written: %d, expected: %d",
				writtenBytes,
				expectedWrittenBytes,
			)
		}
	}
}
