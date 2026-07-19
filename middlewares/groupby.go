package middleware

import "fmt"

func Groupby(key string, data []map[string]interface{}) (map[string][]map[string]interface{}, error) {
	// Create a map to hold the grouped data
	result := make(map[string][]map[string]interface{})

	// Iterate over the provided data
	for _, val := range data {
		// Check if the key exists in the current map
		if v, exists := val[key]; exists {
			// Convert the value to string for grouping
			groupKey := fmt.Sprintf("%v", v) // Convert to string representation
			result[groupKey] = append(result[groupKey], val)
		} else {
			// Handle missing key cases by grouping under an empty string
			result[""] = append(result[""], val)
		}
	}

	// Return the result map
	return result, nil
}
