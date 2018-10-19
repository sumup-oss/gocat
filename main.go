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
	"fmt"

	log "github.com/sumup-oss/go-pkgs/logger"
	"github.com/sumup-oss/go-pkgs/os"

	"github.com/sumup/gocat/cmd"
	"github.com/sumup/gocat/internal/config"
)

func main() {
	osExecutor := &os.RealOsExecutor{}
	configInstance, err := config.NewConfig()
	if err != nil {
		//nolint:errcheck,staticcheck
		fmt.Fprintf(osExecutor.Stderr(), err.Error())
		osExecutor.Exit(1)
	}

	switch configInstance.LogLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	default:
		log.Fatalf("invalid log level %s. Make sure it's upper-case", configInstance.LogLevel)
	}

	logger := log.GetLogger()
	err = cmd.NewRootCmd(osExecutor, logger).Execute()
	if err == nil {
		return
	}

	//nolint:errcheck,staticcheck
	fmt.Fprintf(osExecutor.Stderr(), err.Error())
	osExecutor.Exit(1)
}
