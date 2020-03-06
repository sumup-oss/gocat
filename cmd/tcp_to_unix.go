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
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/spf13/cobra"
	"github.com/sumup-oss/go-pkgs/logger"

	"github.com/sumup-oss/gocat/internal/relay"
)

func NewTCPToUnixCmd(logger logger.Logger) *cobra.Command {
	var tcpToUnixSocketPath string
	var tcpToUnixAddressPath string
	var bufferSize int
	var tcpToUnixHealthCheckInterval time.Duration

	cmdInstance := &cobra.Command{
		Use:   "tcp-to-unix",
		Short: "relay from a TCP source to unix domain socket",
		Long:  `relay from a TCP source to unix domain socket`,
		RunE: func(command *cobra.Command, args []string) error {
			// nolint: gocritic
			if len(tcpToUnixAddressPath) < 0 {
				return stacktrace.NewError("blank/empty `src` specified")
			}

			// nolint: gocritic
			if len(tcpToUnixSocketPath) < 0 {
				return stacktrace.NewError("blank/empty `dst` specified")
			}

			relayer, err := relay.NewTCPtoUnixSocket(
				logger,
				tcpToUnixHealthCheckInterval,
				tcpToUnixAddressPath,
				tcpToUnixSocketPath,
				bufferSize,
			)
			if err != nil {
				return stacktrace.Propagate(err, "couldn't create relay from TCP to unix socket")
			}

			osSignalCh := make(chan os.Signal, 1)
			defer close(osSignalCh)

			signal.Notify(osSignalCh, os.Interrupt, syscall.SIGTERM)

			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			// Ctrl+C handler
			go func() {
				<-osSignalCh
				signal.Stop(osSignalCh)

				_ = os.RemoveAll(tcpToUnixSocketPath)
				cancelFunc()
			}()

			err = relayer.Relay(ctx)
			if err != nil {
				_ = os.RemoveAll(tcpToUnixSocketPath)
				return stacktrace.Propagate(err, "couldn't relay from TCP to unix socket")
			}

			_ = os.RemoveAll(tcpToUnixSocketPath)
			return nil
		},
	}

	cmdInstance.Flags().DurationVar(
		&tcpToUnixHealthCheckInterval,
		"health-check-interval",
		30*time.Second,
		"health check interval for `src`, e.g values are 30m, 60s, 1h.",
	)
	cmdInstance.Flags().StringVar(&tcpToUnixAddressPath, "src", "", "source of TCP address")
	_ = cmdInstance.MarkFlagRequired("src")
	cmdInstance.Flags().StringVar(
		&tcpToUnixSocketPath,
		"dst",
		"",
		"destination of unix domain socket",
	)
	_ = cmdInstance.MarkFlagRequired("dst")
	cmdInstance.Flags().IntVar(
		&bufferSize,
		"buffer-size",
		DefaultBufferSize,
		"Buffer size in bytes of the data stream",
	)
	return cmdInstance
}
