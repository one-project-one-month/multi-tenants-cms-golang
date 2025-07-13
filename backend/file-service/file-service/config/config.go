package config

import (
	utils "github.com/multi-tenants-cms-golang/file-service/pkg"
)

type Configuration struct {
	Port        string
	GitHubToken string
	GitHubOwner string
	GitHubRepo  string
}

func NewConfiguration() *Configuration {
	port := utils.GetEnv("PORT", ":9007")
	githubToken := utils.GetEnv("GITHUB_TOKEN", "")
	githubOwner := utils.GetEnv("GITHUB_OWNER", "")
	githubRepo := utils.GetEnv("GITHUB_REPO", "")
	return &Configuration{
		Port:        port,
		GitHubToken: githubToken,
		GitHubOwner: githubOwner,
		GitHubRepo:  githubRepo,
	}
}
