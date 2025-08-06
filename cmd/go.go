package cmd

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/gddo/gosrc"
	"github.com/spf13/cobra"
	"github.com/xMoelletschi/goreadme/goreadme"
)

var (
	// Holds configuration for Goreadme invocation.
	cfg goreadme.Config

	// Write readme output
	out io.WriteCloser = os.Stdout
)

// goCmd represents the go command
var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Generate a README.md from Go documentation",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client := http.DefaultClient
		gr := goreadme.New(client)

		err := gr.WithConfig(cfg).Create(ctx, pkg(flag.Args()), out)
		if err != nil {
			log.Fatalf("Failed: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(goCmd)

	goCmd.Flags().StringVar(&cfg.ImportPath, "import-path", "", "Override package import path.")
	goCmd.Flags().StringVar(&cfg.Title, "title", "", "Override readme title. Default is package name.")
	goCmd.Flags().StringVar(&cfg.GoDocURL, "godoc-url", "https://pkg.go.dev", "Go Doc URL for GoDoc badge.")
	goCmd.Flags().BoolVar(&cfg.RecursiveSubPackages, "recursive", false, "Load docs recursively.")
	goCmd.Flags().BoolVar(&cfg.RenderTypeContent, "render-type-content", false, "If 'types' is specified, render full type content.")
	goCmd.Flags().BoolVar(&cfg.Consts, "constants", false, "Write package constants section, and if 'types' is specified, also write per-type constants section.")
	goCmd.Flags().BoolVar(&cfg.Vars, "variables", false, "Write package variables section, and if 'types' is specified, also write per-type variables section.")
	goCmd.Flags().BoolVar(&cfg.Functions, "functions", false, "Write functions section.")
	goCmd.Flags().BoolVar(&cfg.Types, "types", false, "Write types section.")
	goCmd.Flags().BoolVar(&cfg.Factories, "factories", false, "If 'types' is specified, write section for functions returning each type.")
	goCmd.Flags().BoolVar(&cfg.Methods, "methods", false, "If 'types' is specified, write section for methods for each type.")
	goCmd.Flags().BoolVar(&cfg.SkipExamples, "skip-examples", false, "Skip the examples section.")
	goCmd.Flags().BoolVar(&cfg.SkipSubPackages, "skip-sub-packages", false, "Skip the sub packages section.")
	goCmd.Flags().BoolVar(&cfg.Badges.TravisCI, "badge-travisci", false, "Show TravisCI badge.")
	goCmd.Flags().BoolVar(&cfg.Badges.CodeCov, "badge-codecov", false, "Show CodeCov badge.")
	goCmd.Flags().BoolVar(&cfg.Badges.GolangCI, "badge-golangci", false, "Show GolangCI badge.")
	goCmd.Flags().BoolVar(&cfg.Badges.GoDoc, "badge-godoc", false, "Show GoDoc badge.")
	goCmd.Flags().BoolVar(&cfg.Badges.GoReportCard, "badge-goreportcard", false, "Show GoReportCard badge.")
	goCmd.Flags().BoolVar(&cfg.GeneratedNotice, "generated-notice", false, "Add generated file notice (visible only in Markdown code).")
	goCmd.Flags().BoolVar(&cfg.Credit, "credit", true, "Add credit line.")
}

func pkg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	path, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}
	gosrc.SetLocalDevMode(path)
	return "."
}
