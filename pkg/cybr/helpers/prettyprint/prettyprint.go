package prettyprint

import (
	"encoding/json"
	"fmt"
)

// PrintJSON will pretty print any data structure to a JSON blob
func PrintJSON(obj interface{}) error {
	json, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}
