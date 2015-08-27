package script

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

// Writes the netrc file.
func writeNetrc(machine, login, password string) string {
	var buf bytes.Buffer
	if len(machine) == 0 {
		return buf.String()
	}
	out := fmt.Sprintf(
		netrcScript,
		machine,
		login,
		password,
	)
	buf.WriteString(out)
	return buf.String()
}

// Writes the RSA private key
func writeKey(key []byte) string {
	var buf bytes.Buffer
	if len(key) == 0 {
		return buf.String()
	}
	keystr := string(key)
	buf.WriteString(fmt.Sprintf(keyScript, keystr))
	buf.WriteString(keyConfScript)
	buf.WriteString(forceYesScript)
	return buf.String()
}

// traceCommand is a helper function that allows us to echo
// commands back to the console for debugging purposes.
func traceCommand(cmd string) string {
	trace := fmt.Sprintf("$ %s\n", cmd)
	encoded := base64.StdEncoding.EncodeToString([]byte(trace))
	return fmt.Sprintf("echo %s | base64 -d\n%s\n", encoded, cmd)
}

// wrapCommand is a helper function that base64 encodes
// a shell command (or entire script)
func wrapCommand(script []byte) string {
	encoded := base64.StdEncoding.EncodeToString(script)
	return fmt.Sprintf("echo %s | base64 -d | $SHELL", encoded)
}
