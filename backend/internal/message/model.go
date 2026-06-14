package message

import "time"

type MessageType string

const (
	TypeText   MessageType = "text"
	TypeSystem MessageType = "system"
)

type Message struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	ChannelID uint        `gorm:"index;not null" json:"channel_id"`
	UserID    uint        `gorm:"not null" json:"user_id"`
	Content   string      `gorm:"not null" json:"content"`
	Type      MessageType `gorm:"default:text" json:"type"`
	IsDeleted bool        `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time   `gorm:"index" json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
