package gatherer

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

func Github(uri string) (string, error) {
	repoURI := ensureRootRepoURL(uri)

	doc, err := htmlquery.LoadURL(repoURI)

	if err != nil {
		return "", fmt.Errorf("unable to load %s: %s", repoURI, err.Error())
	}

	link := htmlquery.FindOne(doc, "//h3[text()='License']/following-sibling::div/a")
	if link == nil {
		return "", fmt.Errorf("unable to determine license from %s", repoURI)
	}

	text := strings.TrimSpace(link.LastChild.Data)
	if strings.Contains(text, "View") {
		href := htmlquery.SelectAttr(link, "href")

		if !strings.HasPrefix(href, "http") {
			href = "https://github.com" + href
		}

		text = href
	}

	return text, nil
}

func ensureRootRepoURL(uri string) string {
	parsed, _ := url.Parse(uri)

	pathParts := strings.Split(parsed.Path, "/")

	if len(pathParts) < 2 {
		return uri
	}

	parsed.Path = strings.Join(pathParts[:3], "/")

	return parsed.String()
}
