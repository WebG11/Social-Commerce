package utils

import (
	"github.com/google/uuid"
)

// 生成 UUID
func GenerateID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
