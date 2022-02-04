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
	"github.com/sumup-oss/go-pkgs/logger"
	"github.com/sumup-oss/go-pkgs/os"

	"github.com/spf13/cobra"
)

const (
	// NOTE: 16k since Linux OS is mostly setting this
	DefaultBufferSize = 16384
)

func NewRootCmd(osExecutor os.OsExecutor, logger logger.Logger) *cobra.Command {
	cmdInstance := &cobra.Command{
		Use:   "gocat",
		Short: "gocat cli utility",
		Long:  "Golang alternative for simple unix pipe to tcp exposure, similar to socat.",
		// NOTE: Silence errors and usage since it'll log twice,
		// due to bad cobra API design and the fact that `RunE` actually returns the error
		// that it's going to log either way.
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmdInstance.AddCommand(
		NewFakeCmd(logger),
		NewTCPToUnixCmd(logger),
		NewTCPToTCPcmd(logger),
		NewUnixToTCPCmd(logger),
		NewVersionCmd(osExecutor),
	)
	return cmdInstance
}
