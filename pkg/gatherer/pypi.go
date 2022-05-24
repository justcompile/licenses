package gatherer

import (
	"context"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
)

type PythonError struct {
	message string
	URL     string
}

func (p *PythonError) Error() string {
	return p.message
}

const pypiURL = "https://pypi.org/project"

func Pypi(ctx context.Context, packageName string) (string, error) {
	packageURI := fmt.Sprintf("%s/%s", pypiURL, packageName)

	doc, err := htmlquery.LoadURL(packageURI)
	if err != nil {
		return "", fmt.Errorf("could not get pypi package information for %s: %s", packageName, err.Error())
	}

	links := htmlquery.Find(doc, "//h3[text()='Project links']/following-sibling::ul/li/a[contains(@href, 'github')]")
	if len(links) == 0 {
		return "", &PythonError{message: fmt.Sprintf("could not location Github link for package %s", packageName), URL: packageURI}
	}

	for _, link := range links {
		if childDataContains(link.LastChild.Data, "Code", "Source", "Home", "Homepage") {
			return htmlquery.SelectAttr(link, "href"), nil
		}
	}

	return "", &PythonError{message: fmt.Sprintf("no github link found %s", packageName), URL: packageURI}
}

func childDataContains(data string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(data, kw) {
			return true
		}
	}

	return false
}
