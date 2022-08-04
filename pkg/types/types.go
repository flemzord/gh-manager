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
