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
	"github.com/yugabyte/yb-voyager/yb-voyager/src/metadb"
	"github.com/yugabyte/yb-voyager/yb-voyager/src/utils"
)

var fallBackSynchronizeCmd = &cobra.Command{
	Use:   "synchronize",
	Short: "This command exports the changes from YugabyteDB.",
	Long:  `This command connects to YugabyteDB and exports the changes received by it so that they can be imported into the fall back database.`,

	Run: func(cmd *cobra.Command, args []string) {
		source.DBType = YUGABYTEDB
		exportType = CHANGES_ONLY
		exporterRole = TARGET_DB_EXPORTER_FB_ROLE
		exportDataCmd.PreRun(cmd, args)
		err := metaDB.UpdateMigrationStatusRecord(func(record *metadb.MigrationStatusRecord) {
			record.FallBackSyncStarted = true
		})
		if err != nil {
			utils.ErrExit("failed to update migration status record for fall-back sync started: %v", err)
		}
		exportDataCmd.Run(cmd, args)
	},
}

func init() {
	fallBackCmd.AddCommand(fallBackSynchronizeCmd)
	registerCommonGlobalFlags(fallBackSynchronizeCmd)
	registerTargetDBAsSourceConnFlags(fallBackSynchronizeCmd)
	registerExportDataFlags(fallBackSynchronizeCmd)
	hideFlagsInFallFowardCmds(fallBackSynchronizeCmd)
	hideExportFlagsInFallBackCmds(fallBackSynchronizeCmd)
}

func hideExportFlagsInFallBackCmds(cmd *cobra.Command) {
	cmd.Flags().Lookup("parallel-jobs").Hidden = true
}
