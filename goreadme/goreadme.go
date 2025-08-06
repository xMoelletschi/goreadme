package goreadme

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang/gddo/doc"
	"github.com/pkg/errors"
	"github.com/xMoelletschi/goreadme/internal/markdown"
	"github.com/xMoelletschi/goreadme/internal/template"
)

// New returns a GoReadme object with a custom client.
// client is an HTTP client used to perform the requests. It can be used
// to authenticate github requests, for example, a github client can be used:
//
//	oauth2.NewClient(ctx, oauth2.StaticTokenSource(
//		&oauth2.Token{AccessToken: "...github access token..."
//	))
func New(c *http.Client) *GoReadme {
	return &GoReadme{client: c}
}

// GoReadme enables getting readme.md text from a go package.
type GoReadme struct {
	client *http.Client
	config Config
}

type Config struct {
	// Override readme title. Default is package name.
	Title string `json:"title"`
	// ImportPath is used to override the import path. For example: github.com/user/project,
	// github.com/user/project/package or github.com/user/project/version.
	ImportPath string `json:"import_path"`
	// GoDocURL is the Go Doc URL used in the GoDoc Badge. Default: https://pkg.go.dev.
	GoDocURL string `json:"godoc_url"`
	// Use the standard library comment parser introduced in Go 1.19 to generate the markdown output.
	StdMarkdown bool `json:"std_markdown"`
	// RenderTypeContent will render fulll type content instead of an ellipsis (`{ ... }`).
	RenderTypeContent bool `json:"render_type_content"`
	// Consts will make constants documentation to be added to the README.
	// If Types is specified, constants for each type will also be added to the README.
	Consts bool `json:"consts"`
	// Vars will make exported variables documentation to be added to the README.
	// If Types is specified, exported variables for each type will also be added to the README.
	Vars bool `json:"vars"`
	// Functions will make functions documentation to be added to the README.
	Functions bool `json:"functions"`
	// Types will make types documentation to be added to the README.
	Types bool `json:"types"`
	// Factories will make functions returning a type to be added to the README, if Types is also specified.
	// Has no effect if Types is not specified.
	Factories bool `json:"factories"`
	// Methods will make the methods for a type to be added to the README, if Types is also specified.
	// Has no effect if Types is not specified.
	Methods bool `json:"methods"`
	// SkipExamples will omit the examples section from the README.
	SkipExamples bool `json:"skip_examples"`
	// SkipSubPackages will omit the sub packages section from the README.
	SkipSubPackages bool `json:"skip_sub_packages"`
	// NoDiffBlocks disables marking code blocks as diffs if they start with minus or plus signes.
	NoDiffBlocks bool `json:"no_diff_blocks"`
	// RecursiveSubPackages will retrieved subpackages information recursively.
	// If false, only one level of subpackages will be retrieved.
	RecursiveSubPackages bool `json:"recursive_sub_packages"`
	Badges               struct {
		TravisCI     bool `json:"travis_ci"`
		CodeCov      bool `json:"code_cov"`
		GolangCI     bool `json:"golang_ci"`
		GoDoc        bool `json:"go_doc"`
		GoReportCard bool `json:"go_report_card"`
	} `json:"badges"`
	// GeneratedFileNotice will add a notice (HTML comment) stating that the README is generated and should probably not be edited.
	GeneratedNotice bool `json:"generated_notice"`
	Credit          bool `json:"credit"`
}

// Create writes the content of readme.md to w, with the default client.
// name should be a Go repository name, such as "github.com/xMoelletschi/goreadme".
func Create(ctx context.Context, name string, w io.Writer) error {
	g := GoReadme{client: http.DefaultClient}
	return g.Create(ctx, name, w)
}

// WithConfig returns a copy of the converter with the given configuration.
func (r GoReadme) WithConfig(cfg Config) *GoReadme {
	r.config = cfg
	return &r
}

// Create writes the content of readme.md to w, with r's HTTP client.
// name should be a Go repository name, such as "github.com/xMoelletschi/goreadme".
func (r *GoReadme) Create(ctx context.Context, name string, w io.Writer) error {
	p, err := r.get(ctx, name)
	if err != nil {
		return err
	}
	return template.Execute(w, p, r.config, markdown.OptNoDiff(r.config.NoDiffBlocks))
}

// pkg contains information about a go package, to be used in the template.
type pkg struct {
	Package     *doc.Package
	SubPackages []subPkg
}

// subPkg is information about sub package, to be used in the template.
type subPkg struct {
	Path    string
	Package *doc.Package
}

func (r *GoReadme) get(ctx context.Context, name string) (*pkg, error) {
	log.Printf("Getting %s", name)
	p, err := docGet(ctx, r.client, name, "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting %s", name)
	}
	sort.Strings(p.Subdirectories)

	// If functions were not requested to be added to the readme, add their
	// examples to the main readme.
	if !r.config.Functions {
		for _, f := range p.Funcs {
			for _, e := range f.Examples {
				if e.Name == "" {
					e.Name = f.Name
				}
				if e.Doc == "" {
					e.Doc = f.Doc
				}
				p.Examples = append(p.Examples, e)
			}
		}
	}

	// If types were not requested to be added to the readme, add their
	// examples to the main readme.
	if !r.config.Types {
		for _, f := range p.Types {
			for _, e := range f.Examples {
				if e.Name == "" {
					e.Name = f.Name
				}
				if e.Doc == "" {
					e.Doc = f.Doc
				}
				p.Examples = append(p.Examples, e)
			}
		}
	}

	if p.IsCmd {
		// TODO: make this better
		p.Name = filepath.Base(name)
		p.Doc = strings.TrimPrefix(p.Doc, "Package main is ")
	}

	if override := r.config.Title; override != "" {
		p.Name = override
	}

	if override := r.config.ImportPath; override != "" {
		p.ImportPath = override
	}

	if override := r.config.GoDocURL; override == "" {
		r.config.GoDocURL = "https://pkg.go.dev"
	}

	pkg := &pkg{
		Package: p,
	}

	if !r.config.SkipSubPackages {
		f := subpackagesFetcher{
			importPath: name,
			client:     r.client,
			recursive:  r.config.RecursiveSubPackages,
		}
		pkg.SubPackages, err = f.Fetch(ctx, p)
		if err != nil {
			return nil, err
		}
	}
	debug(pkg)
	return pkg, nil
}

func debug(p *pkg) {
	if os.Getenv("debug") == "" {
		return
	}

	d, _ := json.MarshalIndent(p, "  ", "  ")
	log.Printf("Package data: %s", string(d))
}
