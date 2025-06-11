package unit

import (
	"testing"

	"github.com/bitdance-panic/gobuy/app/services/user/biz/bll"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	// 测试用例
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "正常登录",
			email:    "test@example.com",
			password: "correct_password",
			wantErr:  false,
		},
		{
			name:     "密码错误",
			email:    "test@example.com",
			password: "wrong_password",
			wantErr:  true,
		},
		{
			name:     "用户不存在",
			email:    "nonexistent@example.com",
			password: "any_password",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行测试
			err := bll.UserLogin(tt.email, tt.password)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUserRegister(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "正常注册",
			email:    "new@example.com",
			password: "valid_password",
			wantErr:  false,
		},
		{
			name:     "邮箱已存在",
			email:    "existing@example.com",
			password: "any_password",
			wantErr:  true,
		},
		{
			name:     "邮箱格式错误",
			email:    "invalid_email",
			password: "any_password",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bll.UserRegister(tt.email, tt.password)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
