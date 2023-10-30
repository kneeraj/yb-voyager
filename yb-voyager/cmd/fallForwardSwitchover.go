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

var fallForwardSwitchoverCmd = &cobra.Command{
	Use:   "switchover",
	Short: "Initiates the switchover to fall-forward DB",
	Long:  `Initiates the switchover to fall-forward DB`,

	Run: func(cmd *cobra.Command, args []string) {
		err := InitiatePrimarySwitch("fallforward")
		if err != nil {
			utils.ErrExit("failed to initiate fallforward: %v", err)
		}
	},
}

func init() {
	fallForwardCmd.AddCommand(fallForwardSwitchoverCmd)
	fallForwardSwitchoverCmd.Flags().StringVarP(&exportDir, "export-dir", "e", "",
		"export directory is the workspace used to keep the exported schema, data, state, and logs")
}
