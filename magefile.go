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

// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

func Test() error {
	return sh.RunV("go", "test", ".")
}

func Bench() error {
	return sh.RunV("go", "test", "-benchtime=5s", "-bench", ".")
}

func BenchAndGraph() error {
	// NOTE: Use https://github.com/miry/benchgraph for this to work.
	return sh.RunV(
		// HACK: Use a shell to perform the piping.
		// Perhaps do it in pure Golang if we need to support non-UNIX.
		"bash",
		"-c",
		`go test -timeout=60m -benchtime=5s -bench=. | benchgraph -title='Benchmark results in ns/op (lower is better)' -function-signature-pattern='Benchmark(?P<functionName>[\w+]+)/(?P<functionArguments>[\w+]+)-(?P<numberOfThreads>[0-9]+)$'`,
	)
}

func BenchTCPToUnixAndGraph() error {
	// NOTE: Use https://github.com/miry/benchgraph for this to work.
	return sh.RunV(
		// HACK: Use a shell to perform the piping.
		// Perhaps do it in pure Golang if we need to support non-UNIX.
		"bash",
		"-c",
		`go test -timeout=60m -benchtime=5s -bench=^BenchmarkTCPToUnix . | benchgraph -title='Benchmark results in ns/op (lower is better)' -function-signature-pattern='Benchmark(?P<functionName>[\w+]+)/(?P<functionArguments>[\w+]+)-(?P<numberOfThreads>[0-9]+)$'`,
	)
}

func BenchUnixToTCPAndGraph() error {
	// NOTE: Use https://github.com/miry/benchgraph for this to work.
	return sh.RunV(
		// HACK: Use a shell to perform the piping.
		// Perhaps do it in pure Golang if we need to support non-UNIX.
		"bash",
		"-c",
		`go test -timeout=60m -benchtime=5s -bench=^BenchmarkUnixToTCP . | benchgraph -title='Benchmark results in ns/op (lower is better)' -function-signature-pattern='Benchmark(?P<functionName>[\w+]+)/(?P<functionArguments>[\w+]+)-(?P<numberOfThreads>[0-9]+)$'`,
	)
}
