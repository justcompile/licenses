package gatherer

import "context"

func NoOp(ctx context.Context, packageName string) (string, error) {
	return "", nil
}
