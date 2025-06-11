package dao

import (
	"gorm.io/gorm"
)

func GetOrderColumns(db *gorm.DB) ([]string, error) {
	// 定义结构体接收列信息
	type ColumnInfo struct {
		Field string `gorm:"column:Field"`
	}
	var columnInfos []ColumnInfo
	// 执行原生查询并扫描到结构体切片
	if err := db.Raw("SHOW COLUMNS FROM `order`").Scan(&columnInfos).Error; err != nil {
		return nil, err
	}
	// 提取列名到字符串切片
	columns := make([]string, len(columnInfos))
	for i, info := range columnInfos {
		columns[i] = info.Field
	}
	return columns, nil
}

func GetProductColumns(db *gorm.DB) ([]string, error) {
	// 定义结构体接收列信息
	type ColumnInfo struct {
		Field string `gorm:"column:Field"`
	}
	var columnInfos []ColumnInfo
	// 执行原生查询并扫描到结构体切片
	if err := db.Raw("SHOW COLUMNS FROM `product`").Scan(&columnInfos).Error; err != nil {
		return nil, err
	}
	// 提取列名到字符串切片
	columns := make([]string, len(columnInfos))
	for i, info := range columnInfos {
		columns[i] = info.Field
	}
	return columns, nil
}
