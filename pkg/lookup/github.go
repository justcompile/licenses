package lookup

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
)

func Github(url string) (string, error) {
	doc, err := htmlquery.LoadURL(url)

	if err != nil {
		return "", fmt.Errorf("unable to load %s: %s", url, err.Error())
	}

	link := htmlquery.FindOne(doc, "//h3[text()='License']/following-sibling::div/a")
	if link == nil {
		return "", fmt.Errorf("unable to determine license from %s", url)
	}

	return strings.TrimSpace(link.LastChild.Data), nil
}
