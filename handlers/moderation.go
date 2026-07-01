package handlers

import "strings"

// Список запрещённых слов согласно ТЗ
var blacklist = []string{"qwerty", "йцукен", "zxvbnm"}

// IsAllowed проверяет текст на наличие запрещённых слов (регистронезависимо)
func IsAllowed(text string) bool {
	lower := strings.ToLower(text)
	for _, word := range blacklist {
		if strings.Contains(lower, strings.ToLower(word)) {
			return false
		}
	}
	return true
}
