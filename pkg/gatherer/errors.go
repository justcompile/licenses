package gatherer

import "fmt"

type RepoNotFound struct {
	message string
	URL     string
}

func (r *RepoNotFound) Error() string {
	return r.message
}

type LicenseFound struct {
	License string
}

func (l *LicenseFound) Error() string {
	return fmt.Sprintf("license found: %s", l.License)
}
