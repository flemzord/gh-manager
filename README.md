# Github Manager

## How to install
### macOS
```
brew install flemzord/tap/gh-manager
```
### Linux & Windows
Download binaries from [Release page](https://github.com/flemzord/gh-manager/releases).

## Usage
Add github_team_name with permission to all repositories in your github_organization.
```
gh-manager repo permissions add --token github_token --organization github_organization --teamName github_team_name --permission admin
```