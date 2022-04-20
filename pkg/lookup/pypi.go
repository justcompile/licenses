package lookup

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
)

const pypiURL = "https://pypi.org/project"

func Pypi(packageName string) (string, error) {
	doc, err := htmlquery.LoadURL(fmt.Sprintf("%s/%s", pypiURL, packageName))
	if err != nil {
		return "", fmt.Errorf("could not get pypi package information for %s: %s", packageName, err.Error())
	}

	links := htmlquery.Find(doc, "//h3[text()='Project links']/following-sibling::ul/li/a[contains(@href, 'github')]")
	if len(links) == 0 {
		return "", fmt.Errorf("could not location Github link for package %s", packageName)
	}

	for _, link := range links {
		if strings.Contains(link.LastChild.Data, "Source") || strings.Contains(link.LastChild.Data, "Homepage") {
			return htmlquery.SelectAttr(link, "href"), nil
		}
	}

	return "", fmt.Errorf("no github link found %s", packageName)
}
