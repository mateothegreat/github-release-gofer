package config

import (
	"regexp"

	"github.com/google/go-github/v63/github"
	"github.com/mateothegreat/go-config/config"
)

// Config is the configuration for the tailer.
type Config struct {
	// Path is the default path to install the files to for all repositories.
	Path string `yaml:"path"`
	// Mode is the default mode to install the file(s) with for all repositories.
	Mode int `yaml:"mode" required:"false"`
	// Repos is a list of repositories to download releases from.
	Repos []Repo `yaml:"repos"`
}

// Repo is a repository to download releases from.
type Repo struct {
	// Name is the name of this repo config.
	Name string `yaml:"name"`
	// Owner is the owner of this repo config.
	Owner string `yaml:"owner"`
	// Repo is the repo of the repository to download from.
	Repo string `yaml:"repo"`
	// Path is the path to install the files to for this repository.
	Path string `yaml:"path"`
	// Download is the pattern to match against the release assets.
	Download string `yaml:"download"`
	// Mode is the mode to install the file(s) with for this repository.
	Mode int `yaml:"mode" required:"false"`
	// Matchers is the pattern to match against the extracted files to install.
	Matchers []Matcher `yaml:"matchers"`
}

// Matcher is a pattern to match against the extracted files to install.
type Matcher struct {
	// Pattern is the pattern to match against the extracted files to install.
	Pattern string `yaml:"pattern"`
	// Mode is the mode to install the file in.
	Mode int `yaml:"mode" required:"false"`
	// Path is the path to install the file to.
	Path string `yaml:"path"`
}

// Match finds the assets that matches the pattern.
//
// Arguments:
//   - assets []*github.ReleaseAsset: The assets to match against the pattern.
//
// Returns:
//   - A slice of assets that match the pattern.
func (r *Repo) Match(assets []*github.ReleaseAsset) []*github.ReleaseAsset {
	matches := []*github.ReleaseAsset{}
	for _, asset := range assets {
		matched, err := regexp.MatchString(r.Download, asset.GetName())
		if err != nil {
			continue
		}
		if matched {
			matches = append(matches, asset)
		}
	}
	return matches
}

// Match checks if the given file matches the pattern.
//
// Arguments:
//   - file string: The file to match against the pattern.
//
// Returns:
//   - A boolean indicating if the file matches the pattern.
func (m *Matcher) Match(file string) bool {
	matched, err := regexp.MatchString(m.Pattern, file)
	if err != nil {
		return false
	}
	return matched
}

// GetConfig returns a config of type T.
// It will merge the base config with the environment config.
// If the environment config does not exist, it will use the base config.
//
// Arguments:
//   - env: The environment to use.
//
// Returns:
//   - A pointer to the config of type T.
//   - An error if the config could not be found.
func GetConfig() (*Config, error) {
	config, err := config.GetConfig[Config](config.GetConfigArgs{
		Paths: []string{
			".gofer.yaml",
			".github-gofer.yaml",
		},
		WalkDepth: 6,
	})
	if err != nil {
		return nil, err
	}
	return config, nil
}
