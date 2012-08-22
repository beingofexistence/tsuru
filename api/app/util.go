package app

import (
	"bytes"
	"fmt"
	"regexp"
)

// filterOutput filters output from juju.
//
// It removes all lines that does not represent useful output, like juju's
// logging and Python's deprecation warnings.
func filterOutput(output []byte) []byte {
	var result [][]byte
	var ignore bool
	deprecation := []byte("DeprecationWarning")
	regexLog := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`)
	regexSshWarning := regexp.MustCompile(`^Warning: Permanently added`)
	lines := bytes.Split(output, []byte{'\n'})
	for _, line := range lines {
		if ignore {
			ignore = false
			continue
		}
		if bytes.Contains(line, deprecation) {
			ignore = true
			continue
		}
		if !regexSshWarning.Match(line) && !regexLog.Match(line) {
			result = append(result, line)
		}
	}
	return bytes.Join(result, []byte{'\n'})
}

// newUUID generates an uuid.
func newUUID() (string, error) {
	f, err := filesystem().Open("/dev/urandom")
	if err != nil {
		return "", err
	}
	b := make([]byte, 16)
	_, err = f.Read(b)
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x", b)
	return uuid, nil
}
