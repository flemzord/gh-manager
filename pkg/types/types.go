package types

import (
	"context"

	"github.com/google/go-github/v45/github"
)

type Config struct {
	Context      context.Context
	Client       github.Client
	Organization string
	Token        string
}

type Permissions []struct {
	Team       string `yaml:"team"`
	Permission string `yaml:"permission"`
	TeamID     int64  `yaml:"team_id,omitempty"`
}
