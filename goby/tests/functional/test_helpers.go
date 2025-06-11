package functional

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// 创建测试请求
func createTestRequest(method, path string, payload interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// 执行测试请求
func executeTestRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// 解析响应
func parseResponse(rr *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	return response, err
}

// 创建测试用户
func createTestUser(email, password string) error {
	// 实现创建测试用户的逻辑
	return nil
}

// 清理测试数据
func cleanupTestData() error {
	// 实现清理测试数据的逻辑
	return nil
}
