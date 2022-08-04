package cmd

import (
	"github.com/spf13/cobra"
)

var organization string
var token string

// repoCmd represents the env command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage your Github Repository",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.PersistentFlags().StringVarP(&organization, "organization", "o", "", "Github Organization")
	repoCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Github Token")
	repoCmd.MarkPersistentFlagRequired("organization")
	repoCmd.MarkPersistentFlagRequired("token")
}
