package functional

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginAPI(t *testing.T) {
	// 测试用例
	tests := []struct {
		name       string
		payload    map[string]string
		wantStatus int
		wantErr    bool
	}{
		{
			name: "正常登录",
			payload: map[string]string{
				"email":    "test@example.com",
				"password": "correct_password",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "密码错误",
			payload: map[string]string{
				"email":    "test@example.com",
				"password": "wrong_password",
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "缺少参数",
			payload: map[string]string{
				"email": "test@example.com",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			jsonData, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(LoginHandler)

			// 执行请求
			handler.ServeHTTP(rr, req)

			// 检查状态码
			assert.Equal(t, tt.wantStatus, rr.Code)

			// 检查响应内容
			if !tt.wantErr {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
			}
		})
	}
}

func TestRegisterAPI(t *testing.T) {
	tests := []struct {
		name       string
		payload    map[string]string
		wantStatus int
		wantErr    bool
	}{
		{
			name: "正常注册",
			payload: map[string]string{
				"email":    "new@example.com",
				"password": "valid_password",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "邮箱已存在",
			payload: map[string]string{
				"email":    "existing@example.com",
				"password": "any_password",
			},
			wantStatus: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "邮箱格式错误",
			payload: map[string]string{
				"email":    "invalid_email",
				"password": "any_password",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(RegisterHandler)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)

			if !tt.wantErr {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			}
		})
	}
}
