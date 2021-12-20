package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var genericPlaceholder, _ = regexp.Compile(`(?mU)<(.*)>`)

var specificPathPlaceholder, _ = regexp.Compile(`(?mU)<path:([^#]+)#([^#]+)(?:#([^#]+))?>`)

// var indivPlaceholderSyntax, _ = regexp.Compile(`(?mU)path:(?P<path>[^#]+?)#(?P<key>[^#]+?)(?:#(?P<version>.+?))??`)

func replaceInner(node *map[string]interface{}, replacerFunc func(string, string) (interface{}, []error)) {
	obj := *node

	for key, value := range obj {
		valueType := reflect.ValueOf(value).Kind()
		// Recurse through nested maps
		if valueType == reflect.Map {
			inner, ok := value.(map[string]interface{})
			if !ok {
				continue
			}
			replaceInner(&inner, replacerFunc)
		} else if valueType == reflect.String {
			replacement, err := replacerFunc(key, value.(string))
			fmt.Println("======")
			fmt.Println("replacement: ", replacement)
			if err != nil {
				panic(err)
			}
			obj[key] = replacement
		}
	}
}

func configReplacement(key, value string) (interface{}, []error) {
	res, err := genericReplacement(key, value)
	if err != nil {
		return nil, err
	}

	// configMap data values must be strings
	return stringify(res), err
}

func genericReplacement(key, value string) (_ interface{}, err []error) {
	var nonStringReplacement interface{}
	var placeholderRegex = genericPlaceholder

	res := placeholderRegex.ReplaceAllFunc([]byte(value), func(match []byte) []byte {
		placeholder := strings.Trim(string(match), "<>")

		// Split modifiers from placeholder
		pipelineFields := strings.Split(placeholder, "|")
		placeholder = strings.Trim(pipelineFields[0], " ")

		var secretValue interface{}

		secretValue = "verysecrettext"

		if secretValue != nil {
			switch secretValue.(type) {
			case string:
				{
					fmt.Println("hereee")
					return []byte(secretValue.(string))
				}
			default:
				{
					nonStringReplacement = secretValue
					return match
				}
			}
		} else {
			missingKeyErr := errors.New("missing secret value")
			err = append(err, missingKeyErr)
		}

		return match
	})

	// The above block can only replace <placeholder> strings with other strings
	// In the case where the value is a non-string, we insert it directly here.
	// Useful for cases like `replicas: <replicas>`
	if nonStringReplacement != nil {
		return nonStringReplacement, err
	}

	return string(res), err
}

func stringify(input interface{}) string {
	switch input.(type) {
	case int:
		{
			return strconv.Itoa(input.(int))
		}
	case bool:
		{
			return strconv.FormatBool(input.(bool))
		}
	case json.Number:
		{
			return string(input.(json.Number))
		}
	case []byte:
		{
			return string(input.([]byte))
		}
	default:
		{
			return input.(string)
		}
	}
}
