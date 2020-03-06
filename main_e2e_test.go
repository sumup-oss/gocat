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

package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sumup-oss/go-pkgs/os"
	"github.com/sumup-oss/go-pkgs/task"
	"github.com/sumup-oss/go-pkgs/testutils"
	gocatTesting "github.com/sumup-oss/gocat/internal/testing"
	"io/ioutil"
	"net"
	stdOs "os"
	"os/exec"
	"testing"
	"time"
)

var (
	osExecutor      = &os.RealOsExecutor{}
	gocatBinaryPath string
)

func hasSocatBinary() error {
	_, err := exec.LookPath("socat")
	return err
}

func TestMain(m *testing.M) {
	gocatBinaryPath = testutils.GoBuild(
		context.Background(),
		"gocat",
		"github.com/sumup-oss/gocat",
		osExecutor,
	)
	runTests := m.Run()

	_ = osExecutor.Remove(gocatBinaryPath)
	stdOs.Exit(runTests)
}

func TestGocatTCPToUnix(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	payload := "123456"
	payloadLength := len(payload)
	dstClient := prepareGocatTCPToUnixTest(ctx, t, payloadLength)
	defer dstClient.Close()

	var sendBuffer bytes.Buffer
	_, err := sendBuffer.Write([]byte(payload))
	require.Nil(t, err, "Failed to write testcase payload in buffer")

	sentPayload := sendBuffer.Bytes()
	n, err := dstClient.SendMsg(sentPayload)
	require.Nil(t, err, "Failed to send payload to gocat dst address")
	require.Equal(
		t,
		payloadLength,
		n,
		"Failed to send complete payload to gocat dst address",
	)

	receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
	require.Nil(t, err, "Failed to receive payload from gocat dst address")
	require.Equal(
		t,
		sentPayload,
		receivedPayload,
		"Different sent compared to received payload",
	)
}

func TestGocatUnixToTCP(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	payload := "123456"
	payloadLength := len(payload)
	dstClient := prepareGocatUnixToTCPTest(ctx, t, payloadLength)
	defer dstClient.Close()

	var sendBuffer bytes.Buffer
	_, err := sendBuffer.Write([]byte(payload))
	require.Nil(t, err, "Failed to write testcase payload in buffer")

	sentPayload := sendBuffer.Bytes()
	n, err := dstClient.SendMsg(sentPayload)
	require.Nil(t, err, "Failed to send payload to gocat dst address")
	require.Equal(
		t,
		payloadLength,
		n,
		"Failed to send complete payload to gocat dst address",
	)

	receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
	require.Nil(t, err, "Failed to receive payload from gocat dst address")
	require.Equal(
		t,
		sentPayload,
		receivedPayload,
		"Different sent compared to received payload",
	)
}

func BenchmarkTCPToUnixSequential_Socat(b *testing.B) {
	err := hasSocatBinary()
	if err != nil {
		b.Skip("Skipping socat benchmark since no socat is present in $PATH")
	}

	testCases := []struct {
		payload string
	}{
		{testutils.RandString(1)},
		{testutils.RandString(10)},
		{testutils.RandString(100)},
		{testutils.RandString(1000)},
		{testutils.RandString(1500)},
		{testutils.RandString(3000)},
		{testutils.RandString(6000)},
		{testutils.RandString(9000)},
		{testutils.RandString(18000)},
		{testutils.RandString(36000)},
		{testutils.RandString(64000)},
	}

	for _, testCase := range testCases {
		ctx, cancelCtx := context.WithCancel(context.Background())

		payloadLength := len(testCase.payload)
		dstClient := prepareSocatTCPToUnixTest(ctx, b, payloadLength)
		// NOTE: Intentional defer in loop to make sure after tests are run/cancelled it's still called
		defer dstClient.Close()

		var sendBuffer bytes.Buffer
		_, err := sendBuffer.Write([]byte(testCase.payload))
		require.Nil(b, err, "Failed to write testcase payload in buffer")

		b.Run(
			fmt.Sprintf("%d", payloadLength),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					sentPayload := sendBuffer.Bytes()
					n, err := dstClient.SendMsg(sentPayload)
					require.Equal(
						b,
						payloadLength,
						n,
						"Failed to send complete payload to gocat dst address",
					)
					require.Nil(b, err, "Failed to send payload to gocat dst address")

					receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
					require.Nil(
						b,
						err,
						"Failed to receive payload from gocat dst address",
					)

					require.Equal(
						b,
						sentPayload,
						receivedPayload,
						"Different sent compared to received payload",
					)
				}
			},
		)

		dstClient.Close()
		cancelCtx()
	}
}

func BenchmarkTCPToUnixSequential_Gocat(b *testing.B) {
	testCases := []struct {
		payload string
	}{
		{testutils.RandString(1)},
		{testutils.RandString(10)},
		{testutils.RandString(100)},
		{testutils.RandString(1000)},
		{testutils.RandString(1500)},
		{testutils.RandString(3000)},
		{testutils.RandString(6000)},
		{testutils.RandString(9000)},
		{testutils.RandString(18000)},
		{testutils.RandString(36000)},
		{testutils.RandString(64000)},
	}

	for _, testCase := range testCases {
		ctx, cancelCtx := context.WithCancel(context.Background())

		payloadLength := len(testCase.payload)
		dstClient := prepareGocatTCPToUnixTest(ctx, b, payloadLength)
		// NOTE: Intentional defer in loop to make sure after tests are run/cancelled it's still called
		defer dstClient.Close()

		var sendBuffer bytes.Buffer
		_, err := sendBuffer.Write([]byte(testCase.payload))
		require.Nil(b, err, "Failed to write testcase payload in buffer")

		b.Run(
			fmt.Sprintf("%d", payloadLength),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					sentPayload := sendBuffer.Bytes()
					n, err := dstClient.SendMsg(sentPayload)
					require.Equal(
						b,
						payloadLength,
						n,
						"Failed to send complete payload to gocat dst address",
					)
					require.Nil(b, err, "Failed to send payload to gocat dst address")

					receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
					require.Nil(
						b,
						err,
						"Failed to receive payload from gocat dst address",
					)

					require.Equal(
						b,
						sentPayload,
						receivedPayload,
						"Different sent compared to received payload",
					)
				}
			},
		)

		dstClient.Close()
		cancelCtx()
	}
}

func BenchmarkTCPToUnixParallel_Gocat(b *testing.B) {
	testCases := []struct {
		payload string
	}{
		{testutils.RandString(1)},
		{testutils.RandString(10)},
		{testutils.RandString(100)},
		{testutils.RandString(1000)},
		{testutils.RandString(1500)},
		{testutils.RandString(3000)},
		{testutils.RandString(6000)},
		{testutils.RandString(9000)},
		{testutils.RandString(18000)},
		{testutils.RandString(36000)},
		{testutils.RandString(64000)},
	}

	for _, testCase := range testCases {
		ctx, cancelCtx := context.WithCancel(context.Background())

		payloadLength := len(testCase.payload)
		dstClient := prepareGocatTCPToUnixTest(ctx, b, payloadLength)
		// NOTE: Intentional defer in loop to make sure after tests are run/cancelled it's still called
		defer dstClient.Close()

		var sendBuffer bytes.Buffer
		_, err := sendBuffer.Write([]byte(testCase.payload))
		require.Nil(
			b,
			err,
			"Failed to write testcase payload in buffer",
		)

		b.Run(
			fmt.Sprintf("%d", payloadLength),
			func(b *testing.B) {
				sentPayloads := make([][]byte, b.N)
				receivedPayloads := make([][]byte, b.N)

				b.ResetTimer()

				go func() {
					for i := 0; i < b.N; i++ {
						sentPayload := sendBuffer.Bytes()
						n, err := dstClient.SendMsg(sentPayload)
						require.Equal(
							b,
							payloadLength,
							n,
							"Failed to send complete payload to gocat dst"+
								" address",
						)

						sentPayloads = append(sentPayloads, sentPayload)
						require.Nil(
							b,
							err,
							"Failed to send payload to gocat dst address",
						)
					}
				}()

				for i := 0; i < b.N; i++ {
					receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
					require.Nil(b, err, "Failed to receive payload")

					receivedPayloads = append(receivedPayloads, receivedPayload)
				}

				require.Equal(b, len(sentPayloads), len(receivedPayloads))

				for i := 0; i < b.N; i++ {
					assert.Equal(b, sentPayloads[i], receivedPayloads[i])
				}
			},
		)

		dstClient.Close()
		cancelCtx()
	}
}

func BenchmarkTCPToUnixParallel_Socat(b *testing.B) {
	testCases := []struct {
		payload string
	}{
		{testutils.RandString(1)},
		{testutils.RandString(10)},
		{testutils.RandString(100)},
		{testutils.RandString(1000)},
		{testutils.RandString(1500)},
		{testutils.RandString(3000)},
		{testutils.RandString(6000)},
		{testutils.RandString(9000)},
		{testutils.RandString(18000)},
		{testutils.RandString(36000)},
		{testutils.RandString(64000)},
	}

	for _, testCase := range testCases {
		ctx, cancelCtx := context.WithCancel(context.Background())

		payloadLength := len(testCase.payload)
		dstClient := prepareSocatTCPToUnixTest(ctx, b, payloadLength)
		// NOTE: Intentional defer in loop to make sure after tests are run/cancelled it's still called
		defer dstClient.Close()

		var sendBuffer bytes.Buffer
		_, err := sendBuffer.Write([]byte(testCase.payload))
		require.Nil(
			b,
			err,
			"Failed to write testcase payload in buffer",
		)

		b.Run(
			fmt.Sprintf("%d", payloadLength),
			func(b *testing.B) {
				sentPayloads := make([][]byte, b.N)
				receivedPayloads := make([][]byte, b.N)

				b.ResetTimer()

				go func() {
					for i := 0; i < b.N; i++ {
						sentPayload := sendBuffer.Bytes()
						n, err := dstClient.SendMsg(sentPayload)
						require.Equal(
							b,
							payloadLength,
							n,
							"Failed to send complete payload to gocat dst"+
								" address",
						)

						sentPayloads = append(sentPayloads, sentPayload)
						require.Nil(
							b,
							err,
							"Failed to send payload to gocat dst address",
						)
					}
				}()

				for i := 0; i < b.N; i++ {
					receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
					require.Nil(b, err, "Failed to receive payload")

					receivedPayloads = append(receivedPayloads, receivedPayload)
				}

				require.Equal(b, len(sentPayloads), len(receivedPayloads))

				for i := 0; i < b.N; i++ {
					assert.Equal(b, sentPayloads[i], receivedPayloads[i])
				}
			},
		)

		dstClient.Close()
		cancelCtx()
	}
}

func BenchmarkUnixToTCPSequential_Gocat(b *testing.B) {
	testCases := []struct {
		payload string
	}{
		{testutils.RandString(1)},
		{testutils.RandString(10)},
		{testutils.RandString(100)},
		{testutils.RandString(1000)},
		{testutils.RandString(1500)},
		{testutils.RandString(3000)},
		{testutils.RandString(6000)},
		{testutils.RandString(9000)},
		{testutils.RandString(18000)},
		{testutils.RandString(36000)},
		{testutils.RandString(64000)},
	}

	for _, testCase := range testCases {
		ctx, cancelCtx := context.WithCancel(context.Background())

		payloadLength := len(testCase.payload)
		dstClient := prepareGocatUnixToTCPTest(ctx, b, payloadLength)
		// NOTE: Intentional defer in loop to make sure after tests are run/cancelled it's still called
		defer dstClient.Close()

		var sendBuffer bytes.Buffer
		_, err := sendBuffer.Write([]byte(testCase.payload))
		require.Nil(b, err, "Failed to write testcase payload in buffer")

		b.Run(
			fmt.Sprintf("%d", payloadLength),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					sentPayload := sendBuffer.Bytes()
					n, err := dstClient.SendMsg(sentPayload)
					require.Equal(
						b,
						payloadLength,
						n,
						"Failed to send complete payload to gocat dst address",
					)
					require.Nil(b, err, "Failed to send payload to gocat dst address")

					receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
					require.Nil(
						b,
						err,
						"Failed to receive payload from gocat dst address",
					)

					require.Equal(
						b,
						sentPayload,
						receivedPayload,
						"Different sent compared to received payload",
					)
				}
			},
		)

		dstClient.Close()
		cancelCtx()
	}
}

func BenchmarkUnixToTCPSequential_Socat(b *testing.B) {
	testCases := []struct {
		payload string
	}{
		{testutils.RandString(1)},
		{testutils.RandString(10)},
		{testutils.RandString(100)},
		{testutils.RandString(1000)},
		{testutils.RandString(1500)},
		{testutils.RandString(3000)},
		{testutils.RandString(6000)},
		{testutils.RandString(9000)},
		{testutils.RandString(18000)},
		{testutils.RandString(36000)},
		{testutils.RandString(64000)},
	}

	for _, testCase := range testCases {
		ctx, cancelCtx := context.WithCancel(context.Background())

		payloadLength := len(testCase.payload)
		dstClient := prepareSocatUnixToTCPTest(ctx, b, payloadLength)
		// NOTE: Intentional defer in loop to make sure after tests are run/cancelled it's still called
		defer dstClient.Close()

		var sendBuffer bytes.Buffer
		_, err := sendBuffer.Write([]byte(testCase.payload))
		require.Nil(b, err, "Failed to write testcase payload in buffer")

		b.Run(
			fmt.Sprintf("%d", payloadLength),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					sentPayload := sendBuffer.Bytes()
					n, err := dstClient.SendMsg(sentPayload)
					require.Equal(
						b,
						payloadLength,
						n,
						"Failed to send complete payload to gocat dst address",
					)
					require.Nil(b, err, "Failed to send payload to gocat dst address")

					receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
					require.Nil(
						b,
						err,
						"Failed to receive payload from gocat dst address",
					)

					require.Equal(
						b,
						sentPayload,
						receivedPayload,
						"Different sent compared to received payload",
					)
				}
			},
		)

		dstClient.Close()
		cancelCtx()
	}
}

func BenchmarkUnixToTCPParallel_Gocat(b *testing.B) {
	testCases := []struct {
		payload string
	}{
		{testutils.RandString(1)},
		{testutils.RandString(10)},
		{testutils.RandString(100)},
		{testutils.RandString(1000)},
		{testutils.RandString(1500)},
		{testutils.RandString(3000)},
		{testutils.RandString(6000)},
		{testutils.RandString(9000)},
		{testutils.RandString(18000)},
		{testutils.RandString(36000)},
		{testutils.RandString(64000)},
	}

	for _, testCase := range testCases {
		ctx, cancelCtx := context.WithCancel(context.Background())

		payloadLength := len(testCase.payload)
		dstClient := prepareGocatUnixToTCPTest(ctx, b, payloadLength)
		// NOTE: Intentional defer in loop to make sure after tests are run/cancelled it's still called
		defer dstClient.Close()

		var sendBuffer bytes.Buffer
		_, err := sendBuffer.Write([]byte(testCase.payload))
		require.Nil(
			b,
			err,
			"Failed to write testcase payload in buffer",
		)

		b.Run(
			fmt.Sprintf("%d", payloadLength),
			func(b *testing.B) {
				sentPayloads := make([][]byte, b.N)
				receivedPayloads := make([][]byte, b.N)

				b.ResetTimer()

				go func() {
					for i := 0; i < b.N; i++ {
						sentPayload := sendBuffer.Bytes()
						n, err := dstClient.SendMsg(sentPayload)
						require.Equal(
							b,
							payloadLength,
							n,
							"Failed to send complete payload to gocat dst"+
								" address",
						)

						sentPayloads = append(sentPayloads, sentPayload)
						require.Nil(
							b,
							err,
							"Failed to send payload to gocat dst address",
						)
					}
				}()

				for i := 0; i < b.N; i++ {
					receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
					require.Nil(b, err, "Failed to receive payload")

					receivedPayloads = append(receivedPayloads, receivedPayload)
				}

				require.Equal(b, len(sentPayloads), len(receivedPayloads))

				for i := 0; i < b.N; i++ {
					assert.Equal(b, sentPayloads[i], receivedPayloads[i])
				}
			},
		)

		dstClient.Close()
		cancelCtx()
	}
}

func BenchmarkUnixToTCPParallel_Socat(b *testing.B) {
	testCases := []struct {
		payload string
	}{
		{testutils.RandString(1)},
		{testutils.RandString(10)},
		{testutils.RandString(100)},
		{testutils.RandString(1000)},
		{testutils.RandString(1500)},
		{testutils.RandString(3000)},
		{testutils.RandString(6000)},
		{testutils.RandString(9000)},
		{testutils.RandString(18000)},
		{testutils.RandString(36000)},
		{testutils.RandString(64000)},
	}

	for _, testCase := range testCases {
		ctx, cancelCtx := context.WithCancel(context.Background())

		payloadLength := len(testCase.payload)
		dstClient := prepareSocatUnixToTCPTest(ctx, b, payloadLength)
		// NOTE: Intentional defer in loop to make sure after tests are run/cancelled it's still called
		defer dstClient.Close()

		var sendBuffer bytes.Buffer
		_, err := sendBuffer.Write([]byte(testCase.payload))
		require.Nil(
			b,
			err,
			"Failed to write testcase payload in buffer",
		)

		b.Run(
			fmt.Sprintf("%d", payloadLength),
			func(b *testing.B) {
				sentPayloads := make([][]byte, b.N)
				receivedPayloads := make([][]byte, b.N)

				b.ResetTimer()

				go func() {
					for i := 0; i < b.N; i++ {
						sentPayload := sendBuffer.Bytes()
						n, err := dstClient.SendMsg(sentPayload)
						require.Equal(
							b,
							payloadLength,
							n,
							"Failed to send complete payload to gocat dst"+
								" address",
						)

						sentPayloads = append(sentPayloads, sentPayload)
						require.Nil(
							b,
							err,
							"Failed to send payload to gocat dst address",
						)
					}
				}()

				for i := 0; i < b.N; i++ {
					receivedPayload, err := dstClient.ReceiveMsg(payloadLength)
					require.Nil(b, err, "Failed to receive payload")

					receivedPayloads = append(receivedPayloads, receivedPayload)
				}

				require.Equal(b, len(sentPayloads), len(receivedPayloads))

				for i := 0; i < b.N; i++ {
					assert.Equal(b, sentPayloads[i], receivedPayloads[i])
				}
			},
		)

		dstClient.Close()
		cancelCtx()
	}
}

func prepareGocatTCPToUnixTest(
	ctx context.Context,
	t gocatTesting.TestingT,
	bufferSize int,
) *gocatTesting.UnixSocketClient {
	binaryBuild := testutils.NewBuild(gocatBinaryPath, "")

	fd, err := ioutil.TempFile("", "gocat-tcp-to-unix-test")
	require.Nil(t, err, "Failed to create temporary file")

	dstListenAddress := fd.Name()

	err = stdOs.RemoveAll(fd.Name())
	require.Nil(t, err, "Failed to delete temporary file")

	testSrcServer := gocatTesting.NewTCPServer(t, bufferSize, "127.0.0.1:0")
	serverListenCh := make(chan *gocatTesting.ListenResult, 1)
	go testSrcServer.Serve(serverListenCh)
	testSrcServerListenResult := <-serverListenCh
	require.Nil(t, testSrcServerListenResult.Err, "Failed to listen with TCP src server")

	go func() {
		stdout, stderr, err := binaryBuild.Run(
			ctx,
			"tcp-to-unix",
			"--src",
			testSrcServerListenResult.Address,
			"--dst",
			dstListenAddress,
		)
		if err != nil {
			fmt.Printf(
				"Failed to run TCP to unix command, stdout: %s, stderr: %s, err: %s\n",
				stdout,
				stderr,
				err,
			)
		}
	}()

	var dstClient *gocatTesting.UnixSocketClient

	// NOTE: Wait for UNIX server to be brought up by gocat
	currentRetries := 0
	clientFn := task.Retry(1*time.Second, func(ctx context.Context) error {
		currentRetries += 1

		dstClient, err = gocatTesting.NewUnixClient(dstListenAddress)
		if err != nil {
			if currentRetries <= 30 {
				return task.NewRetryableError(err)
			}

			return err
		}

		return nil
	})

	err = clientFn(ctx)
	require.Nil(t, err)

	return dstClient
}

func prepareSocatTCPToUnixTest(
	ctx context.Context,
	t gocatTesting.TestingT,
	bufferSize int,
) *gocatTesting.UnixSocketClient {
	fd, err := ioutil.TempFile("", "socat-tcp-to-unix-test")
	require.Nil(t, err, "Failed to create temporary file")

	dstListenAddress := fd.Name()

	err = stdOs.RemoveAll(fd.Name())
	require.Nil(t, err, "Failed to delete temporary file")

	testSrcServer := gocatTesting.NewTCPServer(t, bufferSize, "127.0.0.1:0")
	serverListenCh := make(chan *gocatTesting.ListenResult, 1)
	go testSrcServer.Serve(serverListenCh)
	testSrcServerListenResult := <-serverListenCh
	require.Nil(t, testSrcServerListenResult.Err, "Failed to listen with TCP src server")

	go func() {
		var stdout, stderr bytes.Buffer
		socatCmd := exec.CommandContext(
			ctx,
			"socat",
			"-d",
			fmt.Sprintf(
				"UNIX-LISTEN:%s,unlink-early,mode=777,fork,sndbuf-late=16384,rcvbuf-late=16384",
				dstListenAddress,
			),
			fmt.Sprintf("TCP:%s", testSrcServerListenResult.Address),
		)
		socatCmd.Stdout = &stdout
		socatCmd.Stderr = &stderr

		err := socatCmd.Run()
		if err != nil {
			fmt.Printf(
				"Failed to run TCP to unix command, stdout: %s, stderr: %s, err: %s\n",
				stdout.String(),
				stderr.String(),
				err,
			)
		}
	}()

	var dstClient *gocatTesting.UnixSocketClient

	// NOTE: Wait for UNIX server to be brought up by gocat
	currentRetries := 0
	clientFn := task.Retry(1*time.Second, func(ctx context.Context) error {
		currentRetries += 1

		dstClient, err = gocatTesting.NewUnixClient(dstListenAddress)
		if err != nil {
			if currentRetries <= 30 {
				return task.NewRetryableError(err)
			}

			return err
		}

		return nil
	})

	err = clientFn(ctx)
	require.Nil(t, err)

	return dstClient
}

func prepareSocatUnixToTCPTest(
	ctx context.Context,
	t gocatTesting.TestingT,
	bufferSize int,
) *gocatTesting.TCPClient {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err, "Failed to create temporary address")
	dstListenAddress := l.Addr().String()
	dstListenHost, dstListenPort, err := net.SplitHostPort(dstListenAddress)
	require.Nil(t, err, "Failed to parse and split host from port of address")

	err = l.Close()
	require.Nil(t, err, "Failed to close temporary TCP listener")

	testSrcServer := gocatTesting.NewUnixServer(t, bufferSize)
	srcServerListenCh := make(chan *gocatTesting.ListenResult, 1)
	go testSrcServer.Serve(srcServerListenCh)

	testSrcServerListenResult := <-srcServerListenCh
	require.Nil(t, testSrcServerListenResult.Err, "Failed to listen with Unix socket src server")

	go func() {
		var stdout, stderr bytes.Buffer
		socatCmd := exec.CommandContext(
			ctx,
			"socat",
			"-d",
			fmt.Sprintf(
				"TCP-LISTEN:%s,reuseaddr,fork,range=%s/32,sndbuf-late=16384,rcvbuf-late=16384",
				dstListenPort,
				dstListenHost,
			),
			fmt.Sprintf("UNIX-CLIENT:%s", testSrcServerListenResult.Address),
		)
		socatCmd.Stdout = &stdout
		socatCmd.Stderr = &stderr

		err := socatCmd.Run()
		if err != nil {
			fmt.Printf(
				"Failed to run unix to TCP command, stdout: %s, stderr: %s, err: %s\n",
				stdout.String(),
				stderr.String(),
				err,
			)
		}
	}()

	var dstClient *gocatTesting.TCPClient

	// NOTE: Wait for UNIX server to be brought up by gocat
	currentRetries := 0
	clientFn := task.Retry(1*time.Second, func(ctx context.Context) error {
		currentRetries += 1

		dstClient, err = gocatTesting.NewTCPClient(dstListenAddress)
		if err != nil {
			if currentRetries <= 30 {
				return task.NewRetryableError(err)
			}

			return err
		}

		return nil
	})

	err = clientFn(ctx)
	require.Nil(t, err)

	return dstClient
}

func prepareGocatUnixToTCPTest(
	ctx context.Context,
	t gocatTesting.TestingT,
	bufferSize int,
) *gocatTesting.TCPClient {
	binaryBuild := testutils.NewBuild(gocatBinaryPath, "")

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err, "Failed to create temporary address")
	dstListenAddress := l.Addr().String()

	testSrcServer := gocatTesting.NewUnixServer(t, bufferSize)
	srcServerListenCh := make(chan *gocatTesting.ListenResult, 1)
	go testSrcServer.Serve(srcServerListenCh)

	testSrcServerListenResult := <-srcServerListenCh
	require.Nil(t, testSrcServerListenResult.Err, "Failed to listen with Unix socket src server")

	go func() {
		err = l.Close()
		require.Nil(t, err, "Failed to close temporary TCP listener")
		stdout, stderr, err := binaryBuild.Run(
			ctx,
			"unix-to-tcp",
			"--src",
			testSrcServerListenResult.Address,
			"--dst",
			dstListenAddress,
		)
		if err != nil {
			fmt.Printf(
				"Failed to run TCP to unix command, stdout: %s, stderr: %s, err: %s\n",
				stdout,
				stderr,
				err,
			)
		}
	}()

	var dstClient *gocatTesting.TCPClient

	// NOTE: Wait for UNIX server to be brought up by gocat
	currentRetries := 0
	clientFn := task.Retry(1*time.Second, func(ctx context.Context) error {
		currentRetries += 1

		dstClient, err = gocatTesting.NewTCPClient(dstListenAddress)
		if err != nil {
			if currentRetries <= 30 {
				return task.NewRetryableError(err)
			}

			return err
		}

		return nil
	})

	err = clientFn(ctx)
	require.Nil(t, err)

	return dstClient
}
