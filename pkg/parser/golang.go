package parser

import (
	"bytes"
	"context"
	"io"

	"golang.org/x/mod/modfile"
)

func Golang(ctx context.Context, r io.Reader) ([]*Package, error) {
	buf := new(bytes.Buffer)
	if _, readErr := buf.ReadFrom(r); readErr != nil {
		return nil, readErr
	}

	f, err := modfile.Parse("", buf.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	packages := make([]*Package, len(f.Require))

	for i, mod := range f.Require {
		packages[i] = &Package{
			Name:    mod.Mod.Path,
			Version: mod.Mod.Version,
		}
	}

	return packages, nil
}
