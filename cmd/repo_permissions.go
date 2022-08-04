package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var PermissionsValue = []string{
	"read",
	"triage",
	"write",
	"maintain",
	"admin",
}

// repoPermissionsCmd represents the env command
var repoPermissionsCmd = &cobra.Command{
	Use:   "permissions",
	Short: "Manage Permissions for your Github Repository",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !Contains(PermissionsValue, Permission) {
			return fmt.Errorf("Invalid Permission, Possible values: %v", PermissionsValue)
		}
		return nil
	},
}

func Contains(value []string, permission string) bool {
	for _, v := range value {
		if v == permission {
			return true
		}
	}
	return false
}

var TeamName string
var Permission string

func init() {
	repoCmd.AddCommand(repoPermissionsCmd)
	repoPermissionsCmd.PersistentFlags().StringVarP(&TeamName, "teamName", "n", "Admin", "Your GitHub Team Name")
	repoPermissionsCmd.PersistentFlags().StringVarP(&Permission, "permission", "p", "read", "Your GitHub Repository Permission. Possible values: read, triage, write, maintain, admin")
	repoPermissionsCmd.MarkPersistentFlagRequired("teamName")
	repoPermissionsCmd.MarkPersistentFlagRequired("permission")
}
