package cmd

import (
	"github.com/flemzord/gh-manager/pkg/apply"
	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply config file",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			apply.NewApplyCommand(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(applyCmd)
}
