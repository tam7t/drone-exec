package builder

import "github.com/drone/drone-exec/builder/parse"

func Parse(yaml string) (*Build, error) {
	t, err := parse.Parse(yaml)
	if err != nil {
		return nil, err
	}
	return Load(t)
}

func Load(tree *parse.Tree) (*Build, error) {

	return &Build{tree: tree}, nil
}
