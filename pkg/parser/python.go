package parser

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

func Python(ctx context.Context, r io.Reader) ([]*Package, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	doc, err := toml.LoadBytes(data)

	if err == nil {
		fmt.Println("TOML:", doc.String())
		return nil, nil
	}

	return parsePip(ctx, bytes.NewReader(data))
}

func parsePip(ctx context.Context, r io.Reader) ([]*Package, error) {
	scanner := bufio.NewScanner(r)

	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("unable to read input: %s", err.Error())
	}

	packages := make([]*Package, 0)

	for _, line := range lines {
		if strings.HasPrefix(line, "-r") {
			refFile, err := tryOpenReferencedPipFile(ctx, line)
			if err != nil {
				return nil, err
			}

			refPackages, refErr := parsePip(ctx, refFile)
			if refErr != nil {
				return nil, refErr
			}

			packages = append(packages, refPackages...)
			continue
		}

		if p, err := packageFromPipRequirement(line); err == nil {
			packages = append(packages, p)
		} else {
			return nil, err
		}
	}

	return packages, nil
}

func packageFromPipRequirement(line string) (*Package, error) {
	parts := strings.Split(line, "==")

	p := &Package{
		Name: parts[0],
	}

	if len(parts) > 1 {
		p.Version = parts[1]
	}

	return p.Sanitise()
}

func tryOpenReferencedPipFile(ctx context.Context, line string) (io.Reader, error) {
	fileName := strings.TrimLeft(line, "-r")

	cd := ctx.Value(ContextCurrentWorkingDir)

	return os.Open(fmt.Sprintf("%s/%s", cd, strings.TrimSpace(fileName)))
}
