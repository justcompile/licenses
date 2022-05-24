package parser

import "strings"

type currentWorkingDir string

const (
	ContextCurrentWorkingDir currentWorkingDir = "workingdir"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	License string `json:"license"`
}

func (p *Package) Sanitise() (*Package, error) {
	if strings.Contains(p.Name, "[") {
		p.Name = strings.SplitN(p.Name, "[", 2)[0]
	}

	return p, nil
}
