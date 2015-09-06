package yaml

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	// "strconv"
	// "strings"

	"gopkg.in/yaml.v2"
)

// Parse parses a Yaml configuraiton file.
func Parse(in []byte) (*Config, error) {
	c := Config{}
	e := yaml.Unmarshal(in, &c)
	return &c, e
}

// ParseString parses a Yaml configuration file
// in string format.
func ParseString(in string) (*Config, error) {
	return Parse([]byte(in))
}

// ParseChecksum parses a Yaml configuration file
// in string format and verifies the checksum.
func ParseChecksum(in, checksum string) (*Config, bool, error) {
	conf, err := ParseString(in)
	if err != nil {
		return nil, false, err
	}
	return conf, shasum(in, checksum), nil
}

// shasum is a helper function that calculates
// and verifies a file checksum. This supports
// the sha1, sha256 and sha512 values.
func shasum(in, checksum string) bool {
	var hash string
	// var size int64
	// var name string

	switch len(checksum) {
	case 64:
		hash = sha512sum(in)
	case 128:
		hash = sha512sum(in)
	case 40:
		hash = sha512sum(in)
	case 0:
		return true // if no checksum assume valid
	}

	// // the checksum might be split into multiple
	// // sections including the file size and name.
	// switch strings.Count(checksum, " ") {
	// case 1:
	// 	fmt.Sscanf(checksum, "%s %s", &checksum, &name)
	// case 2:
	// 	fmt.Sscanf(checksum, "%s %d %s", &checksum, &size, &name)
	// }

	// var cksum = parts[0]
	// var sizes = parts[1]
	// size, err := strconv.ParseInt(size, 10, 64)
	// if err != nil {
	// 	return checksum == cksum
	// }

	return checksum == hash
}

func sha1sum(in string) string {
	h := sha1.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func sha256sum(in string) string {
	h := sha256.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func sha512sum(in string) string {
	h := sha512.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}
