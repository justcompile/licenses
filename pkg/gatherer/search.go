package gatherer

import "context"

func Search(ctx context.Context, packageName string) (string, error) {
	gitURL, err := Pypi(ctx, packageName)
	if err != nil {
		return "", err
	}

	return Github(gitURL)
}
