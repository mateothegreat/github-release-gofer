package github

import "github.com/google/go-github/v63/github"

func NewClient(token string) *github.Client {
	var client *github.Client
	if token != "" {
		client = github.NewClient(nil).WithAuthToken(token)
	} else {
		client = github.NewClient(nil)
	}
	return client
}
