package prettyprint

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

var colorMap = map[string]string{
	"reset":  "\033[0m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"gray":   "\033[37m",
	"white":  "\033[97m",
}

// PrintJSON will pretty print any data structure to a JSON blob
func PrintJSON(obj interface{}) error {
	json, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}

// PrintColor will print a message in a given color
func PrintColor(color string, message string) error {
	if runtime.GOOS == "windows" {
		fmt.Println(message)
		return nil
	}
	if _, ok := colorMap[strings.ToLower(color)]; !ok {
		return fmt.Errorf("invalid color: %s", color)
	}

	fmt.Println(colorMap[strings.ToLower(color)] + message + colorMap["reset"])

	return nil
}
