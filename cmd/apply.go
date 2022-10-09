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

	repos, err := pkgGithub.GetAllRepository(&githubAuth)
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
			err = pkgGithub.AddTeamPermissionToRepository(&githubAuth, permission[0].Team, *repo.Name, permission[0].Permission)
			if err != nil {
				log.Errorf("Error: %s", err)
				return
			}
		}
		// Clean Teams
		teams, _ := pkgGithub.GetTeamsForRepository(&githubAuth, *repo.Name)
		for _, team := range teams {
			if !Contains(githubTeams, *team.Name) && !Contains(configFile.Repository.ExcludeTeam, *team.Name) {
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
