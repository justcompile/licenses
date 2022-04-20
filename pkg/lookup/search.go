package lookup

func Search(packageName string) (string, error) {
	gitURL, err := Pypi(packageName)
	if err != nil {
		return "", err
	}

	return Github(gitURL)
}
