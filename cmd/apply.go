package cmd

import (
	"context"
	pkgGithub "github.com/flemzord/gh-manager/pkg/github"
	"github.com/flemzord/gh-manager/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply config file",
		Long:  ``,
		Run:   applyConfig,
	}
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

func applyConfig(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	githubAuth := pkgGithub.Login(ctx, viper.GetString("github.organization"), viper.GetString("github.token"))
	// Get permissions and retrieve group id
	var permissions types.Permissions
	err := viper.UnmarshalKey("permissions", &permissions)
	if err != nil {
		return
	}
	var githubTeams []string
	for index, permission := range permissions {
		id, _ := pkgGithub.GetTeamID(&githubAuth, permission.Team)
		permission.TeamID = id
		permissions[index] = permission
		githubTeams = append(githubTeams, permission.Team)
	}
	log.Debugln(permissions)
	log.Debugln(githubTeams)

	repos, err := pkgGithub.GetAllRepository(&githubAuth)
	if err != nil {
		log.Errorf("Error: %s", err)
	}
	for _, repo := range repos {
		log.Infof("Repository: %s", *repo.Name)
		archived := repo.GetArchived()
		if archived {
			continue
		}
		if err != nil {
			log.Errorf("Error: %s", err)
			return
		}
		for _, permission := range permissions {
			err = pkgGithub.AddTeamPermissionToRepository(&githubAuth, permission.Team, permission.TeamID, *repo.Name, permission.Permission)
			if err != nil {
				log.Errorf("Error: %s", err)
				return
			}
		}
		// Clean Teams
		teams, _ := pkgGithub.GetTeamsForRepository(&githubAuth, *repo.Name)
		for _, team := range teams {
			if !Contains(githubTeams, *team.Name) {
				err = pkgGithub.RemoveTeamPermissionToRepository(&githubAuth, *team.Name, *repo.Name)
				if err != nil {
					log.Errorf("Error: %s", err)
					return
				}
			}
		}
		// Remove all solo Collaborators
		collaborators, _ := pkgGithub.GetCollaboratorsForRepository(&githubAuth, *repo.Name)
		for _, collaborator := range collaborators {
			log.Printf("Collaborator %s as unauthorized access to the %s repository\n", *collaborator.Login, *repo.Name)
			err := pkgGithub.RemoveCollaboratorPermissionToRepository(&githubAuth, *repo.Name, *collaborator.Login)
			if err != nil {
				log.Errorf("Error: %s", err)
			}
		}
	}
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
