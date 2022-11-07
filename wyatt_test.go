package wyatt_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yardbirdsax/wyatt"
)

func TestUnmarshal(t *testing.T) {
	t.Run("pointer not passed", func(t *testing.T) {
		obj := struct{}{}

		err := wyatt.Unmarshal(obj)

		assert.Error(t, err)
	})

	t.Run("simple input types", func(t *testing.T) {
		type simpleInputStruct struct {
			StringInput string  `json:"string"`
			BoolInput   bool    `json:"bool"`
			IntInput    int     `json:"int"`
			FloatInput  float64 `json:"float"`
		}

		inputs := map[string]string{
			"STRING": "a string",
			"BOOL":   "true",
			"INT":    "32",
			"FLOAT":  "24.32",
		}
		deferFunc, err := setupEnvVars(inputs)
		require.Nil(t, err)
		defer deferFunc()
		expectedOutput := simpleInputStruct{
			StringInput: "a string",
			BoolInput:   true,
			IntInput:    32,
			FloatInput:  24.32,
		}
		actualOutput := simpleInputStruct{}

		err = wyatt.Unmarshal(&actualOutput)
		require.Nil(t, err)

		assert.EqualValues(t, expectedOutput, actualOutput)
	})

	t.Run("complex types", func(t *testing.T) {
		type nestedStruct struct {
			InputA string `json:"inputa"`
			InputB int    `json:"inputb"`
		}
		type complexInputStruct struct {
			ListStringInput []string               `json:"list_string_input"`
			MapStringInput  map[string]interface{} `json:"map_string_input"`
			StructInput     nestedStruct           `json:"struct_input"`
		}
		inputs := map[string]string{
			"LIST_STRING_INPUT": "[\"hello\", \"there\"]",
			"MAP_STRING_INPUT":  "{\"hello\": \"there\"}",
			"STRUCT_INPUT":      "{\"inputA\": \"hello\",\"inputB\": 3}",
		}
		deferFunc, err := setupEnvVars(inputs)
		require.Nil(t, err)
		defer deferFunc()
		expectedOutuput := complexInputStruct{
			ListStringInput: []string{"hello", "there"},
			MapStringInput:  map[string]interface{}{"hello": "there"},
			StructInput: nestedStruct{
				InputA: "hello",
				InputB: 3,
			},
		}
		actualOutput := complexInputStruct{}

		err = wyatt.Unmarshal(&actualOutput)
		require.Nil(t, err)

		assert.EqualValues(t, expectedOutuput, actualOutput)
	})
}

func setupEnvVars(envMap map[string]string) (func(), error) {
	for k, v := range envMap {
		envVarName := fmt.Sprintf("INPUT_%s", k)

		err := os.Setenv(envVarName, v)
		if err != nil {
			return nil, err
		}
	}

	deferFunc := func() {
		for k := range envMap {
			envVarName := fmt.Sprintf("INPUT_%s", k)
			os.Unsetenv(envVarName)
		}
	}

	return deferFunc, nil
}
