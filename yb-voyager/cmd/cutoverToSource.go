/*
Copyright (c) YugabyteDB, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yugabyte/yb-voyager/yb-voyager/src/utils"
)

var cutoverToSourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Initiate cutover to source DB",
	Long:  `Initiate cutover to source DB`,

	Run: func(cmd *cobra.Command, args []string) {
		err := InitiateCutover("source", false, false)
		if err != nil {
			utils.ErrExit("failed to initiate fallback: %v", err)
		}
	},
}

func init() {
	cutoverToCmd.AddCommand(cutoverToSourceCmd)
	registerExportDirFlag(cutoverToSourceCmd)
}
