package cmd

import (
	"context"

	"github.com/fatih/color"
	pkgGithub "github.com/flemzord/gh-manager/pkg/github"
	"github.com/flemzord/gh-manager/pkg/types"
	"github.com/google/go-github/v45/github"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// repoListCmd represents the env command
var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your Github Repository",
	Long:  ``,
	Run:   repoList,
}

func repoList(cmd *cobra.Command, args []string) {
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

	repos, _ := pkgGithub.GetRepository(&config)

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Private", "Archived", "Description")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, repo := range repos {
		tbl.AddRow(repo.GetName(), repo.GetPrivate(), repo.GetArchived(), repo.GetDescription())
	}

	tbl.Print()
}

func init() {
	repoCmd.AddCommand(repoListCmd)
}
