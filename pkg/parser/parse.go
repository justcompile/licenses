package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pelletier/go-toml"
)

func Parse(r io.Reader) ([]*Package, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	doc, err := toml.LoadBytes(data)

	if err == nil {
		fmt.Println("TOML:", doc.String())
		return nil, nil
	}

	return parsePip(bytes.NewReader(data))
}

func parsePip(r io.Reader) ([]*Package, error) {
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
