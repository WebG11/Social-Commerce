package utils

import (
	"strings"
)

// 移除代码块标记
func CleanBlock(rawSQL string) string {
	// 移除末尾的代码块标记（如果有）
	rawSQL = strings.TrimSuffix(rawSQL, "```")
	// 移除开头的 ```html 标记（如果有）
	cleaned := strings.TrimPrefix(rawSQL, "```html")
	if cleaned != rawSQL {
		return cleaned
	}
	cleaned = strings.TrimPrefix(cleaned, "```sql")
	if cleaned != rawSQL {
		return cleaned
	}
	return cleaned
}
