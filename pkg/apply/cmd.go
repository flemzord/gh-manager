package apply

import (
	"context"

	"github.com/flemzord/gh-manager/pkg/github"
	"github.com/flemzord/gh-manager/pkg/lib"
	"github.com/flemzord/gh-manager/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewApplyCommand(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	client := github.Login(ctx, viper.GetString("github.token"))
	githubAuth := types.GithubConfig{
		Context:      ctx,
		Client:       client,
		Organization: viper.GetString("github.organization"),
	}
	// Create struct with Config File
	var configFile types.Config
	err := viper.Unmarshal(&configFile)
	if err != nil {
		return
	}

	var githubTeams []string
	for _, permission := range configFile.Repository.Permissions {
		githubTeams = append(githubTeams, permission[0].Team)
	}
	log.Debugln(configFile)
	log.Debugln(githubTeams)

	repos, err := github.GetAllRepository(&githubAuth)
	if err != nil {
		log.Errorf("Error: %s", err)
	}
	log.Infof("Found %d repositories", len(repos))
	for index, repo := range repos {
		log.Infof("Repository %s updated (%d/%d)", *repo.Name, index+1, len(repos))
		archived := repo.GetArchived()
		if archived {
			continue
		}
		if err != nil {
			log.Errorf("Error: %s", err)
			return
		}
		for _, permission := range configFile.Repository.Permissions {
			err = github.AddTeamPermissionToRepository(&githubAuth, permission[0].Team, *repo.Name, permission[0].Permission)
			if err != nil {
				log.Errorf("Error: %s", err)
				return
			}
		}
		// Clean Teams
		teams, _ := github.GetTeamsForRepository(&githubAuth, *repo.Name)
		for _, team := range teams {
			if !lib.Contains(githubTeams, *team.Name) && !lib.Contains(configFile.Repository.ExcludeTeam, *team.Name) {
				err = github.RemoveTeamPermissionToRepository(&githubAuth, *team.Name, *repo.Name)
				if err != nil {
					log.Errorf("Error: %s", err)
					return
				}
			}
		}
		// Remove all solo Collaborators
		collaborators, _ := github.GetCollaboratorsForRepository(&githubAuth, *repo.Name)
		for _, collaborator := range collaborators {
			log.Printf("Collaborator %s as unauthorized access to the %s repository\n", *collaborator.Login, *repo.Name)
			err := github.RemoveCollaboratorPermissionToRepository(&githubAuth, *repo.Name, *collaborator.Login)
			if err != nil {
				log.Errorf("Error: %s", err)
			}
		}
	}
}
