package utils

import (
	"encoding/json"
	"strings"
)

// IsJSON checks if byte slice is valid JSON
// Returns true if data can be unmarshaled as JSON, false otherwise
func IsJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}

// PrettyJSON formats JSON data with indentation for better readability
// Returns formatted JSON bytes or error if data isn't valid JSON
func PrettyJSON(data []byte) ([]byte, error) {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	
	return json.MarshalIndent(obj, "", "  ")
}

// ParseHeaders converts string array of "key:value" pairs to a map
// Used for processing header flags from command line
func ParseHeaders(headers []string) map[string]string {
	h := make(map[string]string)
	
	for _, header := range headers {
		split := strings.SplitN(header, ":", 2)
		if len(split) == 2 {
			h[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
		}
	}
	
	return h
} 