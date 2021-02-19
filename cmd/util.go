package cmd

import (
	"fmt"
	"strings"
)

// The content will look like
// port=something, sp
func keyValueStringToMap(content string) (map[string]string, error) {
	if content == "" {
		return nil, nil
	}

	if !strings.Contains(content, "=") {
		return nil, fmt.Errorf("Invalid platform prop content. The provided content does not container a '='")
	}

	m := make(map[string]string)

	// TODO: Gotta be a better way to do this
	replaceWith := "^||||^"

	// If the address or property contains a `\,` then replace
	content = strings.ReplaceAll(content, "\\,", replaceWith)
	props := strings.Split(content, ",")
	for _, prop := range props {
		if !strings.Contains(prop, "=") {
			return nil, fmt.Errorf("Property '%s' is invalid because it does not contain a '=' to seperate key from value", prop)
		}
		kvs := strings.SplitN(prop, "=", 2)
		key := strings.Trim(kvs[0], " ")
		value := strings.Trim(strings.ReplaceAll(kvs[1], replaceWith, ","), " ")
		m[key] = value
	}

	return m, nil
}
