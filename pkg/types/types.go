package types

import (
	"context"

	"github.com/google/go-github/v47/github"
)

type GithubConfig struct {
	Context      context.Context
	Client       github.Client
	Organization string
	Token        string
}

type Config struct {
	Dryrun     bool       `yaml:"dryrun,omitempty"`
	Github     Github     `yaml:"github"`
	Repository Repository `yaml:"repository,omitempty"`
}
type Github struct {
	Organization string `yaml:"organization"`
	Token        string `yaml:"token"`
}

type Repository struct {
	Permissions []Permission `yaml:"permissions,omitempty"`
	ExcludeTeam []string     `yaml:"excludeTeam,omitempty"`
}

type Permission []struct {
	Team       string `yaml:"team,omitempty"`
	Permission string `yaml:"permission,omitempty"`
}
