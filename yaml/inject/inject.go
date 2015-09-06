package inject

import (
	"sort"
	"strings"
)

// Inject injects a map of parameters into a raw string and returns
// the resulting string.
//
// Parameters are represented in the string using $$ notation, similar
// to how environment variables are defined in Makefiles.
func Inject(raw string, params map[string]string) string {
	if params == nil {
		return raw
	}
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	injected := raw
	for _, k := range keys {
		v := params[k]
		injected = strings.Replace(injected, "$$"+k, v, -1)
	}
	return injected
}
