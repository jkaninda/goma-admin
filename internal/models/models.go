package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringArray handles text[] PostgreSQL type
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}
	return "{" + joinStrings(a, ",") + "}", nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	str := string(bytes)
	if str == "{}" || str == "" {
		*a = []string{}
		return nil
	}

	// Remove curly braces
	str = str[1 : len(str)-1]
	*a = splitStrings(str, ",")
	return nil
}

// IntArray handles integer[] PostgreSQL type
type IntArray []int

func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}

	result := "{"
	for i, v := range a {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", v)
	}
	result += "}"
	return result, nil
}

func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	str := string(bytes)
	if str == "{}" || str == "" {
		*a = []int{}
		return nil
	}

	// Remove curly braces
	str = str[1 : len(str)-1]
	parts := splitStrings(str, ",")

	result := make([]int, 0, len(parts))
	for _, part := range parts {
		var num int
		fmt.Sscanf(part, "%d", &num)
		result = append(result, num)
	}

	*a = result
	return nil
}

// JSONB handles jsonb PostgreSQL type
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

func joinStrings(arr []string, sep string) string {
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += sep
		}
		result += `"` + s + `"`
	}
	return result
}

func splitStrings(str, sep string) []string {
	if str == "" {
		return []string{}
	}

	var result []string
	var current string
	inQuote := false

	for i := 0; i < len(str); i++ {
		ch := str[i]

		if ch == '"' {
			inQuote = !inQuote
			continue
		}

		if ch == ',' && !inQuote {
			result = append(result, current)
			current = ""
			continue
		}

		current += string(ch)
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}
