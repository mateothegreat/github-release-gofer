package commands

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/go-github/v63/github"
	"github.com/mateothegreat/github-release-gofer/config"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mateothegreat/go-util/files"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/cobra"
)

type File struct {
	Path      string
	Filename  string
	Directory string
	Content   string
}

var Upgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version of one or more packages.",
	Long:  "Upgrade to the latest version of one or more packages.",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			multilog.Fatal("upgrade", "failed to get token", map[string]interface{}{
				"error": err,
			})
		}

		owner, err := cmd.Flags().GetString("owner")
		if err != nil {
			multilog.Fatal("upgrade", "failed to get owner", map[string]interface{}{
				"error": err,
				"owner": owner,
			})
		}

		repo, err := cmd.Flags().GetString("repo")
		if err != nil {
			multilog.Fatal("upgrade", "failed to get repo", map[string]interface{}{
				"error": err,
				"repo":  repo,
			})
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			multilog.Fatal("create", "failed to get path", map[string]interface{}{
				"error": err,
				"path":  path,
			})
		}

		config, err := config.GetConfig()
		if err != nil {
			multilog.Fatal("upgrade", "failed to get config", map[string]interface{}{
				"error": err,
			})
		}

		var client *github.Client
		if token != "" {
			client = github.NewClient(nil).WithAuthToken(token)
		} else {
			client = github.NewClient(nil)
		}

		for _, repo := range config.Repos {
			latest, _, err := client.Repositories.GetLatestRelease(context.Background(), repo.Owner, repo.Name)
			if err != nil {
				multilog.Fatal("upgrade", "failed to get latest release", map[string]interface{}{
					"error": err,
					"repo":  repo,
				})
			}

			assets := repo.Match(latest.Assets)
			if len(assets) == 0 {
				multilog.Warn("upgrade", "no assets found", map[string]interface{}{
					"repo": repo,
				})
				continue
			}

			// Download the asset.
			for _, asset := range assets {
				client := &http.Client{}
				req, err := http.NewRequest("GET", *asset.BrowserDownloadURL, nil)
				if err != nil {
					multilog.Fatal("upgrade", "failed to create request", map[string]interface{}{
						"error": err,
						"repo":  repo,
					})
				}
				req.Header.Set("Accept", "application/octet-stream")

				resp, err := client.Do(req)
				if err != nil {
					multilog.Fatal("upgrade", "failed to download release asset", map[string]interface{}{
						"error": err,
						"repo":  repo,
					})
				}

				// Unarchive the file using archiver
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					multilog.Fatal("upgrade", "failed to read response body", map[string]interface{}{
						"error": err,
						"repo":  repo,
					})
				}
				resp.Body.Close()

				format, _, err := archiver.Identify(*asset.BrowserDownloadURL, bytes.NewReader(body))
				if err != nil {
					multilog.Fatal("upgrade", "failed to identify archive format", map[string]interface{}{
						"error": err,
						"repo":  repo,
					})
				}
				if ex, ok := format.(archiver.Extractor); ok {
					err = ex.Extract(context.Background(), bytes.NewReader(body), nil, func(ctx context.Context, f archiver.File) error {
						for _, matcher := range repo.Matchers {
							if matcher.Match(f.Name()) {
								destPath := filepath.Join(matcher.Path, f.Name())
								mode, err := strconv.ParseUint(matcher.Mode, 8, 32)
								if err != nil {
									return err
								}
								destFile, err := os.OpenFile(files.ExpandPath(destPath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(mode))
								if err != nil {
									return err
								}
								defer destFile.Close()

								rc, err := f.Open()
								if err != nil {
									return err
								}
								defer rc.Close()

								_, err = io.Copy(destFile, rc)
								if err != nil {
									return err
								}
							}
						}
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
	Upgrade.Flags().StringP("token", "t", "", "GitHub token if you have private repos (optional).")
	Upgrade.Flags().StringP("owner", "o", "", "GitHub owner (username or organization).")
	Upgrade.Flags().StringP("repo", "r", "", "GitHub repository name.")
	Upgrade.Flags().StringP("path", "p", "~/.bin", "Path to extract the release to.")
}
