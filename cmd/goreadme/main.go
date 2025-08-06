// Goreadme command line tool and Github action
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/gddo/gosrc"
	"github.com/posener/goaction/log"
	"github.com/xMoelletschi/goreadme"
)

var (
	// Holds configuration for Goreadme invocation.
	cfg goreadme.Config

	// Write readme output
	out io.WriteCloser = os.Stdout
)

func init() {
	flag.StringVar(&cfg.ImportPath, "import-path", "", "Override package import path.")
	flag.StringVar(&cfg.Title, "title", "", "Override readme title. Default is package name.")
	flag.StringVar(&cfg.GoDocURL, "godoc-url", "https://pkg.go.dev", "Go Doc URL for GoDoc badge.")
	flag.BoolVar(&cfg.RecursiveSubPackages, "recursive", false, "Load docs recursively.")
	flag.BoolVar(&cfg.RenderTypeContent, "render-type-content", false, "If 'types' is specified, render full type content.")
	flag.BoolVar(&cfg.Consts, "constants", false, "Write package constants section, and if 'types' is specified, also write per-type constants section.")
	flag.BoolVar(&cfg.Vars, "variables", false, "Write package variables section, and if 'types' is specified, also write per-type variables section.")
	flag.BoolVar(&cfg.Functions, "functions", false, "Write functions section.")
	flag.BoolVar(&cfg.Types, "types", false, "Write types section.")
	flag.BoolVar(&cfg.Factories, "factories", false, "If 'types' is specified, write section for functions returning each type.")
	flag.BoolVar(&cfg.Methods, "methods", false, "If 'types' is specified, write section for methods for each type.")
	flag.BoolVar(&cfg.SkipExamples, "skip-examples", false, "Skip the examples section.")
	flag.BoolVar(&cfg.SkipSubPackages, "skip-sub-packages", false, "Skip the sub packages section.")
	flag.BoolVar(&cfg.Badges.TravisCI, "badge-travisci", false, "Show TravisCI badge.")
	flag.BoolVar(&cfg.Badges.CodeCov, "badge-codecov", false, "Show CodeCov badge.")
	flag.BoolVar(&cfg.Badges.GolangCI, "badge-golangci", false, "Show GolangCI badge.")
	flag.BoolVar(&cfg.Badges.GoDoc, "badge-godoc", false, "Show GoDoc badge.")
	flag.BoolVar(&cfg.Badges.GoReportCard, "badge-goreportcard", false, "Show GoReportCard badge.")
	flag.BoolVar(&cfg.GeneratedNotice, "generated-notice", false, "Add generated file notice (visible only in Markdown code).")
	flag.BoolVar(&cfg.Credit, "credit", true, "Add credit line.")
	flag.Usage = func() {
		fmt.Fprint(
			flag.CommandLine.Output(),
			`goreadme: Create markdown file from go doc.

Usage:
	goreadme [flags] [import path]

import path (optional): Create a readme file for a package from github.
 Omitting import path will create a readme for the package in CWD.
Flags:
`)
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	ctx := context.Background()
	client := http.DefaultClient
	gr := goreadme.New(client)

	err := gr.WithConfig(cfg).Create(ctx, pkg(flag.Args()), out)
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}

	return
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
