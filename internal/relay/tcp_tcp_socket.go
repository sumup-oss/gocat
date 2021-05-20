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

package relay

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/sumup-oss/go-pkgs/logger"
)

type TCPtoTCPsocket struct {
	AbstractDuplexRelay
}

func NewTCPtoTCPsocket(
	logger logger.Logger,
	healthCheckInterval time.Duration,
	srcAddress,
	dstAddress string,
	bufferSize int,
) (*TCPtoTCPsocket, error) {
	tcpAddressParts := strings.Split(srcAddress, ":")
	if len(tcpAddressParts) != 2 {
		return nil, stacktrace.NewError(
			"wrong format for tcp address %s. Expected <addr>:<port>",
			srcAddress,
		)
	}

	_, err := strconv.ParseInt(tcpAddressParts[1], 10, 32)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"could not parse specified port number %s",
			tcpAddressParts[1],
		)
	}

	return &TCPtoTCPsocket{
		AbstractDuplexRelay{
			healthCheckInterval: healthCheckInterval,
			logger:              logger,
			sourceName:          "source TCP connection",
			destinationName:     "destination TCP connection",
			destinationAddr:     dstAddress,
			bufferSize:          bufferSize,
			dialSourceConn: func(ctx context.Context) (net.Conn, error) {
				dialer := &net.Dialer{
					KeepAlive: tcpKeepAlivePeriod,
				}
				conn, err := dialer.DialContext(
					ctx,
					"tcp",
					srcAddress,
				)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to dial TCP address: %s",
						srcAddress,
					)
				}

				tcpConn := conn.(*net.TCPConn)
				// TODO: Re-evaluate if this is redundant when `KeepAlive` and `net.Dialer` is used.
				_ = tcpConn.SetKeepAlive(true)
				_ = tcpConn.SetKeepAlivePeriod(tcpKeepAlivePeriod)
				return tcpConn, nil
			},
			listenTargetConn: func(ctx context.Context) (net.Listener, error) {
				var lc net.ListenConfig
				listener, err := lc.Listen(ctx, "tcp", dstAddress)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to listen at destination TCP address: %s",
						dstAddress,
					)
				}
				return listener, nil
			},
		},
	}, nil
}
