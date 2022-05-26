package examine

import (
	"context"
	"fmt"
	"io"
	"justcompile/licenses/pkg/gatherer"
	"justcompile/licenses/pkg/parser"
	"os"
	"sync"
)

type PackageParser func(context.Context, io.Reader) ([]*parser.Package, error)

type LicenseGatherer func(context.Context, string) (string, error)

type Examiner struct {
	Parser          PackageParser
	LicenseGatherer LicenseGatherer
}

func (e *Examiner) Process(ctx context.Context, r io.Reader) ([]*parser.Package, error) {
	packages, parseErr := e.Parser(ctx, r)

	if parseErr != nil {
		return nil, fmt.Errorf("unable to extract packages: %s", parseErr.Error())
	}

	var wg sync.WaitGroup
	wg.Add(len(packages))

	for _, p := range packages {
		go func(pkg *parser.Package) {
			defer wg.Done()

			license, err := e.getLicense(ctx, pkg.Name)
			if err != nil {
				switch licError := err.(type) {
				case *gatherer.RepoNotFound:
					pkg.License = licError.URL
				case *gatherer.LicenseFound:
					pkg.License = licError.License
					return
				}

				os.Stderr.WriteString(err.Error() + "\n")
				return
			}

			pkg.License = license
		}(p)
	}

	wg.Wait()

	return packages, nil
}

func (e *Examiner) getLicense(ctx context.Context, packageName string) (string, error) {
	githubURI, err := e.LicenseGatherer(ctx, packageName)
	if err != nil {
		return "", err
	}

	return gatherer.Github(githubURI)
}

func New(lang string) (*Examiner, error) {
	var pkgParser PackageParser
	var licenseGatherer LicenseGatherer

	switch lang {
	case "py", "python":
		pkgParser = parser.Python
		licenseGatherer = gatherer.Pypi
	case "go", "golang":
		pkgParser = parser.Golang
		licenseGatherer = gatherer.Golang
	default:
		return nil, fmt.Errorf("%s language is currently unsupported", lang)
	}

	return &Examiner{
		Parser:          pkgParser,
		LicenseGatherer: licenseGatherer,
	}, nil
}
