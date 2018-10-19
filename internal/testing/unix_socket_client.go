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
	"net"
	"time"

	"github.com/sumup/gocat/internal/relay"
)

type UnixSocketClient struct {
	connection net.Conn
}

func NewUnixClient(address string) (*UnixSocketClient, error) {
	c, err := net.Dial("unix", address)
	if err != nil {
		return nil, err
	}

	return &UnixSocketClient{
		connection: relay.NewDeadlineConnection(c, 30*time.Second, 30*time.Second),
	}, nil
}

func (c *UnixSocketClient) SendMsg(msg []byte) (int, error) {
	return c.connection.Write(msg)
}

func (c *UnixSocketClient) Close() {
	_ = c.connection.Close()
}

func (c *UnixSocketClient) ReceiveMsg(bufferSize int) ([]byte, error) {
	offset := 0

	buf := make([]byte, bufferSize, bufferSize+1)

	for offset < len(buf) {
		n, err := c.connection.Read(buf[offset:])
		if err != nil {
			return nil, err
		}

		offset += n
	}

	return buf, nil
}
