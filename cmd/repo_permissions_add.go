package cmd

import (
	"context"

	pkgGithub "github.com/flemzord/gh-manager/pkg/github"
	"github.com/flemzord/gh-manager/pkg/types"
	"github.com/google/go-github/v45/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// repoPermissionsAddCmd represents the env command
var repoPermissionsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add permission for Team to Github Repository",
	Long:  ``,
	Run:   repoAdd,
}

func repoAdd(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	config := types.Config{
		Context:      ctx,
		Client:       *client,
		Organization: organization,
		Token:        token,
	}

	pkgGithub.AddTeamPermissionToRepository(&config, TeamName, Permission)
}

func init() {
	repoPermissionsCmd.AddCommand(repoPermissionsAddCmd)
}
