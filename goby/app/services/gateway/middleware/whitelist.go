package middleware

import (
	"context"
	"regexp"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func WhiteListMiddleware() app.HandlerFunc {
	whiteList := []string{
		"/login",
		"/product/search",
		"/register",
	}
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.URI().Path())
		for _, p := range whiteList {
			if matchPath(p, path) {
				c.Set("skip_auth", true)
				c.Next(ctx) // 跳过后续中间件
				return
			}
		}
		// c.Request.Header.VisitAll(func(key, value []byte) {
		// 	if string(key) == "Authorization" {
		// 		fmt.Println(string(key), string(value))
		// 	}
		// })
		c.Next(ctx)
	}
}

// matchPath 支持通配符匹配（如 /public/*）
func matchPath(pattern, path string) bool {
	// 如果pattern是正则表达式（以^开头）
	if strings.HasPrefix(pattern, "^") {
		matched, _ := regexp.MatchString(pattern, path)
		return matched
	}

	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(path, prefix)
	}
	return path == pattern
}
