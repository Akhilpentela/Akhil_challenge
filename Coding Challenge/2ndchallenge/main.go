package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// JSONData represents the input JSON structure
type JSONData map[string]interface{}

// TransformedData represents the transformed JSON structure
type TransformedData []map[string]interface{}

// transformJSON transforms the input JSON data to the desired output format
func transformJSON(input JSONData) TransformedData {
	var output TransformedData

	transformedMap := make(map[string]interface{})

	for key, value := range input {
		if key == "" {
			continue
		}

		// Transform keys
		key = strings.TrimSpace(key)
		switch key {
		case "number_1":
			transformedMap["number_1"] = transformNumber(value)
		case "string_1":
			transformedMap["string_1"] = transformString(value)
		case "string_2":
			transformedMap["string_2"] = transformString(value)
		case "map_1":
			transformedMap["map_1"] = transformMap(value.(map[string]interface{}))
		}
	}

	output = append(output, transformedMap)

	return output
}

// transformMap transforms a map within the JSON data
func transformMap(m map[string]interface{}) map[string]interface{} {
	transformedMap := make(map[string]interface{})

	for k, v := range m {
		transformedMap[k] = transformValue(k, v)
	}

	return transformedMap
}

// transformValue transforms a value based on its data type
func transformValue(key string, value interface{}) interface{} {
	switch key {
	case "N", "S", "BOOL", "NULL":
		return transformString(value)
	case "L":
		return transformList(value)
	}

	return nil
}

// transformString transforms a string value based on its data type
func transformString(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		for _, s := range v {
			return strings.TrimSpace(fmt.Sprintf("%v", s))
		}
		return nil
	case string:
		// Transform RFC3339 formatted strings to Unix Epoch
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t.Unix()
		}
		return strings.TrimSpace(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// transformNumber transforms a numeric value
func transformNumber(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		// Strip leading zeros
		s := strings.TrimLeft(v, "0")

		// Convert string to float64
		if num, err := strconv.ParseFloat(s, 64); err == nil {
			return num
		}
		return nil
	case float64:
		return v
	default:
		return nil
	}
}

// transformList transforms a list value
func transformList(value interface{}) interface{} {
	var transformedList []interface{}

	switch v := value.(type) {
	case string:
		if v == "noop" || v == "" {
			return nil
		}

		// Parse list string into JSON array
		var listArray []interface{}
		if err := json.Unmarshal([]byte(v), &listArray); err != nil {
			fmt.Println("Error unmarshalling list:", err)
			return nil
		}

		// Transform each item in the list
		for _, item := range listArray {
			switch itemValue := item.(type) {
			case string:
				transformedList = append(transformedList, transformString(itemValue))
			case float64:
				transformedList = append(transformedList, transformNumber(itemValue))
			case bool:
				transformedList = append(transformedList, itemValue)
			}
		}

		return transformedList
	default:
		return nil
	}
}

func main() {
	// Read input JSON from file
	inputJSONFile, err := os.Open("input.json")
	if err != nil {
		fmt.Println("Error opening input JSON file:", err)
		os.Exit(1)
	}
	defer inputJSONFile.Close()

	var inputJSON JSONData
	err = json.NewDecoder(inputJSONFile).Decode(&inputJSON)
	if err != nil {
		fmt.Println("Error decoding input JSON:", err)
		os.Exit(1)
	}

	// Transform JSON data
	outputJSON := transformJSON(inputJSON)

	// Print output to stdout
	output, err := json.MarshalIndent(outputJSON, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling output:", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
