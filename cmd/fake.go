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

package cmd

import (
	"net"
	"os"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/sumup-oss/go-pkgs/logger"

	"github.com/spf13/cobra"
)

func NewFakeCmd(logger logger.Logger) *cobra.Command {
	cmdInstance := &cobra.Command{
		Use:   "fake",
		Short: "fake unix domain socket server",
		Long:  `fake unix domain socket server`,
		RunE: func(command *cobra.Command, args []string) error {
			_ = os.Remove("./test.sock")
			l, err := net.Listen("unix", "./test.sock")
			if err != nil {
				return stacktrace.Propagate(err, "failed to listen to unix domain socket")
			}

			defer os.Remove("./test.sock")
			for {
				conn, err := l.Accept()
				if err != nil {
					logger.Warnf("Connection error: %s", err)
				}

				logger.Infof("Opened connection from remote addr: %s", conn.RemoteAddr())

				go func(innerConn net.Conn) {
					for {
						_, _ = innerConn.Write([]byte("MOSHI MOSHI"))
						time.Sleep(1 * time.Second)
					}
				}(conn)
			}
		},
	}

	return cmdInstance
}
