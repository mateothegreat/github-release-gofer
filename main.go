package main

import (
	"os"

	"github.com/mateothegreat/github-release-gofer/commands"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "github-gofer",
	Short: "GitHub Release Downloading tool.",
	Long:  "Github Gofer is a CLI tool for downloading GitHub Releases.",
}

func init() {
	multilog.RegisterLogger(multilog.LogMethod("console"), multilog.NewConsoleLogger(&multilog.NewConsoleLoggerArgs{
		Level:  multilog.DEBUG,
		Format: multilog.FormatText,
	}))

	root.CompletionOptions.HiddenDefaultCmd = true
	root.PersistentFlags().BoolP("dry-run", "", false, "Dry run the command.")
	root.PersistentFlags().StringP("token", "t", os.Getenv("GITHUB_TOKEN"), "GitHub token.")
}

func main() {
	root.AddCommand(commands.Upgrade)
	root.AddCommand(commands.List)
	root.Execute()
}
