/*
Wyatt is Go library that marshalls / un-marshalls GitHub Action context data into / from structs. It's named for
the (in)famous U.S. Marshall Wyatt Earp: https://history.howstuffworks.com/historical-figures/wyatt-earp.htm.

See the public repository's README file for more examples.
*/
package wyatt

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
UnMarshall populates a tagged Go struct with input data by the GitHub Action's platform.
Tags should be in the form of the lowercased name of the input. For example, an input called
`Hello_There` would need a tag of `json:"hello_there"`. JSON tags are used since the library
converts the inputs to a JSON object, thereby leveraging the existing functionality present
in the standard library's `encoding/json` package.
*/
func UnMarshall(out interface{}) error {
	var (
		err error
	)

	envJSONMap, err := createEnvironmentJSONMap()
	if err != nil {
		return fmt.Errorf("error generataing environment variable map: %w", err)
	}

	err = json.Unmarshal(envJSONMap, out)

	if err != nil {
		return fmt.Errorf("error unmarshalling environment variable map to struct: %w", err)
	}

	return nil
}

func createEnvironmentJSONMap() ([]byte, error) {
	var (
		envMap     = make(map[string]interface{})
		envJSONMap []byte
		err        error
		prefix     = "INPUT_"
	)

	for _, envVariable := range os.Environ() {
		delimIndex := strings.Index(envVariable, "=")
		if delimIndex == -1 {
			continue
		}

		envKey := envVariable[:delimIndex]

		if strings.HasPrefix(envKey, prefix) {
			inputName := envKey[len(prefix):]
			inputValue := envVariable[delimIndex+1:]

			if boolValue, err := strconv.ParseBool(inputValue); err == nil {
				envMap[inputName] = boolValue
				continue
			}

			if floatValue, err := strconv.ParseFloat(inputValue, 64); err == nil {
				envMap[inputName] = floatValue
				continue
			}

			listValue := []interface{}{}

			err := json.Unmarshal([]byte(inputValue), &listValue)
			if err == nil {
				envMap[inputName] = listValue
				continue
			}

			mapValue := map[string]interface{}{}

			err = json.Unmarshal([]byte(inputValue), &mapValue)
			if err == nil {
				envMap[inputName] = mapValue
				continue
			}

			envMap[inputName] = inputValue
		}
	}

	envJSONMap, err = json.Marshal(envMap)

	return envJSONMap, err
}
