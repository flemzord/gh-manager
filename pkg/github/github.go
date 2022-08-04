package github

import (
	"fmt"

	"github.com/flemzord/gh-manager/pkg/types"
	"github.com/google/go-github/v45/github"
	log "github.com/sirupsen/logrus"
)

func GetRepository(config *types.Config) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := config.Client.Repositories.ListByOrg(config.Context, config.Organization, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func GetTeams(config *types.Config) ([]*github.Team, error) {
	opt := &github.ListOptions{}
	teams, _, err := config.Client.Teams.ListTeams(config.Context, config.Organization, opt)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return teams, nil
}

func GetTeamId(config *types.Config, TeamName string) (int64, error) {
	teams, err := GetTeams(config)
	if err != nil {
		return 0, err
	}
	for _, team := range teams {
		if *team.Name == TeamName {
			return *team.ID, nil
		}
	}
	return 0, nil
}

func AddTeamPermissionToRepository(config *types.Config, TeamName string, Permission string) error {
	organizationInfo, _, _ := config.Client.Organizations.Get(config.Context, config.Organization)
	adminTeam, _ := GetTeamId(config, TeamName)
	if adminTeam == 0 {
		return nil
	}
	repos, err := GetRepository(config)
	if err != nil {
		log.Errorf("Error: %s", err)
	}
	for _, repo := range repos {
		fmt.Printf("Team %s has been added to %s with permission %s\n", TeamName, *repo.Name, Permission)

		opt := &github.TeamAddTeamRepoOptions{Permission: Permission}
		_, err := config.Client.Teams.AddTeamRepoByID(config.Context, organizationInfo.GetID(), adminTeam, config.Organization, *repo.Name, opt)
		if err != nil {
			return err
		}
	}
	return nil
}
