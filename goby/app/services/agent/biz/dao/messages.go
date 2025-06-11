package dao

import (
	"github.com/bitdance-panic/gobuy/app/models"
	"gorm.io/gorm"
	"time"
)

type ChatMessage struct {
	ID         int64     `gorm:"primaryKey"`
	SenderID   int64     `gorm:"not null"`
	ReceiverID int64     `gorm:"not null"`
	Content    string    `gorm:"not null"`
	Timestamp  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}

// SendMessage 添加聊天消息到数据库
func SendMessage(db *gorm.DB, senderID, receiverID int64, content string) error {
	message := models.ChatMessage{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Timestamp:  time.Now(),
	}
	return db.Create(&message).Error
}

// GetHistoryMessages 获取历史聊天记录
func GetHistoryMessages(db *gorm.DB, senderID, receiverID int64, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := db.Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).
		Or("sender_id = ? AND receiver_id = ?", receiverID, senderID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}
