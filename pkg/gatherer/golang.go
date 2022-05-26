package gatherer

import (
	"context"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
)

const goPkgURI = "https://pkg.go.dev/"

func Golang(ctx context.Context, packageName string) (string, error) {
	packageURI := fmt.Sprintf("%s/%s", goPkgURI, packageName)

	doc, err := htmlquery.LoadURL(packageURI)
	if err != nil {
		return "", fmt.Errorf("could not get pypi package information for %s: %s", packageName, err.Error())
	}

	licenseLink := htmlquery.FindOne(doc, "//a[@data-test-id='UnitHeader-license']")
	if licenseLink != nil {
		return licenseLink.LastChild.Data, &LicenseFound{License: licenseLink.LastChild.Data}
	}

	if strings.Contains(packageName, "github") {
		return "https://" + packageName, nil
	}

	return "", &RepoNotFound{message: fmt.Sprintf("no github link found %s", packageName), URL: packageURI}
}
