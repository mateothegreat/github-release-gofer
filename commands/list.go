package commands

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v63/github"
	"github.com/mateothegreat/github-release-gofer/config"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/cobra"
)

var List = &cobra.Command{
	Use:   "list",
	Short: "List all the files in a release.",
	Long:  "List all the files in a release.",
	Run: func(cmd *cobra.Command, args []string) {
		// t, err := cmd.Flags().GetString("token")
		// if err != nil {
		// 	multilog.Fatal("upgrade", "failed to get token", map[string]interface{}{
		// 		"error": err,
		// 	})
		// }

		// owner, err := cmd.Flags().GetString("owner")
		// if err != nil {
		// 	multilog.Fatal("upgrade", "failed to get owner", map[string]interface{}{
		// 		"error": err,
		// 		"owner": owner,
		// 	})
		// }

		// repo, err := cmd.Flags().GetString("repo")
		// if err != nil {
		// 	multilog.Fatal("upgrade", "failed to get repo", map[string]interface{}{
		// 		"error": err,
		// 		"repo":  repo,
		// 	})
		// }

		// path, err := cmd.Flags().GetString("path")
		// if err != nil {
		// 	multilog.Fatal("create", "failed to get path", map[string]interface{}{
		// 		"error": err,
		// 		"path":  path,
		// 	})
		// }

		config, err := config.GetConfig()
		if err != nil {
			multilog.Fatal("upgrade", "failed to get config", map[string]interface{}{
				"error": err,
			})
		}

		client := github.NewClient(nil)

		// Get the latest release for each repository.
		for _, repo := range config.Repos {
			// Get the latest release.
			latest, _, err := client.Repositories.GetLatestRelease(context.Background(), repo.Owner, repo.Repo)
			if err != nil {
				multilog.Fatal("list", "failed to get latest release", map[string]interface{}{
					"error": err,
					"repo":  repo,
				})
			}

			// Download the asset.
			for _, asset := range latest.Assets {
				client := &http.Client{}
				req, err := http.NewRequest("GET", *asset.BrowserDownloadURL, nil)
				if err != nil {
					multilog.Fatal("list", "failed to create request", map[string]interface{}{
						"error": err,
						"repo":  repo,
					})
				}
				req.Header.Set("Accept", "application/octet-stream")

				// Download the asset.
				resp, err := client.Do(req)
				if err != nil {
					multilog.Fatal("list", "failed to download release asset", map[string]interface{}{
						"error": err,
						"repo":  repo,
					})
				}

				// Unarchive the file using archiver
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					multilog.Fatal("list", "failed to read response body", map[string]interface{}{
						"error": err,
						"repo":  repo,
					})
				}
				resp.Body.Close()

				format, _, err := archiver.Identify(*asset.BrowserDownloadURL, bytes.NewReader(body))
				if err != nil {
					multilog.Fatal("list", "failed to identify archive format", map[string]interface{}{
						"error": err,
						"repo":  repo,
					})
				}
				if ex, ok := format.(archiver.Extractor); ok {
					err = ex.Extract(context.Background(), bytes.NewReader(body), nil, func(ctx context.Context, f archiver.File) error {
						multilog.Info("list", "found file", map[string]interface{}{
							"file": f.Name(),
							"repo": repo,
						})
						return nil
					})
					if err != nil {
						multilog.Fatal("upgrade", "failed to extract archive", map[string]interface{}{
							"error": err,
							"repo":  repo,
						})
					}
				} else {
					multilog.Fatal("upgrade", "unsupported archive format", map[string]interface{}{
						"repo": repo,
					})
				}

				multilog.Debug("upgrade", "downloaded release asset", map[string]interface{}{
					"url":  *asset.BrowserDownloadURL,
					"repo": repo,
				})
			}
		}
	},
}

func init() {
	// Upgrade.Flags().StringP("token", "t", "", "GitHub token if you have private repos (optional).")
	// Upgrade.Flags().StringP("owner", "o", "", "GitHub owner (username or organization).")
	// Upgrade.Flags().StringP("repo", "r", "", "GitHub repository name.")
	// Upgrade.Flags().StringP("path", "p", "~/.bin", "Path to extract the release to.")
}
