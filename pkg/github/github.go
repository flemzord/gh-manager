package github

import (
	"context"

	"github.com/flemzord/gh-manager/pkg/types"
	"github.com/google/go-github/v47/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func Login(ctx context.Context, organization string, token string) types.GithubConfig {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	config := types.GithubConfig{
		Context:      ctx,
		Client:       *client,
		Organization: organization,
		Token:        token,
	}
	return config
}

func GetAllRepository(config *types.GithubConfig) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
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

func GetTeamsForRepository(config *types.GithubConfig, repoName string) ([]*github.Team, error) {
	teams, _, err := config.Client.Repositories.ListTeams(config.Context, config.Organization, repoName, &github.ListOptions{})
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func GetTeams(config *types.GithubConfig) ([]*github.Team, error) {
	opt := &github.ListOptions{}
	teams, _, err := config.Client.Teams.ListTeams(config.Context, config.Organization, opt)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return teams, nil
}

func GetTeamID(config *types.GithubConfig, TeamName string) (int64, error) {
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

func AddTeamPermissionToRepository(config *types.GithubConfig, TeamName string, Repository string, Permission string) error {
	organizationInfo, _, _ := config.Client.Organizations.Get(config.Context, config.Organization)
	TeamID, _ := GetTeamID(config, TeamName)

	opt := &github.TeamAddTeamRepoOptions{Permission: Permission}
	_, err := config.Client.Teams.AddTeamRepoByID(config.Context, organizationInfo.GetID(), TeamID, config.Organization, Repository, opt)
	if err != nil {
		return err
	}
	log.Debugf("Repository %s: Team %s has been added or updated with permission %s\n", Repository, TeamName, Permission)
	return nil
}

func RemoveTeamPermissionToRepository(config *types.GithubConfig, RepositoryName string, TeamName string) error {
	organizationInfo, _, _ := config.Client.Organizations.Get(config.Context, config.Organization)
	Team, _ := GetTeamID(config, TeamName)

	_, err := config.Client.Teams.RemoveTeamRepoByID(config.Context, organizationInfo.GetID(), Team, config.Organization, RepositoryName)
	if err != nil {
		return err
	}
	log.Debugf("Repository %s: Team %s has been removed\n", RepositoryName, TeamName)
	return nil
}

func RemoveCollaboratorPermissionToRepository(config *types.GithubConfig, RepositoryName string, CollaborateurName string) error {
	_, err := config.Client.Repositories.RemoveCollaborator(config.Context, config.Organization, RepositoryName, CollaborateurName)
	if err != nil {
		return err
	}
	log.Debugf("Repository %s: Collaborator %s has been removed\n", RepositoryName, CollaborateurName)
	return nil
}

func GetCollaboratorsForRepository(config *types.GithubConfig, RepositoryName string) ([]*github.User, error) {
	opt := &github.ListCollaboratorsOptions{
		Affiliation: "direct",
	}
	collaborators, _, err := config.Client.Repositories.ListCollaborators(config.Context, config.Organization, RepositoryName, opt)
	if err != nil {
		return nil, err
	}
	return collaborators, nil
}
