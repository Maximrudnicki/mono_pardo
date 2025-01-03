package utils

import "mono_pardo/pkg/data/request"

func ConvertFieldUpdatesToMap(updates []request.FieldUpdate) map[string]interface{} {
	result := make(map[string]interface{})
	for _, update := range updates {
		result[update.Field] = update.Value
	}
	return result
}
