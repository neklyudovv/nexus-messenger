package channel

import "time"

type ChannelType string

const (
	TypePublic ChannelType = "public"
	TypeDM     ChannelType = "dm"
)

type Channel struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	WorkspaceID uint        `gorm:"index;not null" json:"workspace_id"`
	Name        string      `gorm:"not null" json:"name"`
	Description string      `json:"description"`
	Type        ChannelType `gorm:"not null" json:"type"`
	CreatedBy   uint        `json:"created_by"`
	CreatedAt   time.Time   `json:"created_at"`
}

type ChannelMember struct {
	ChannelID uint      `gorm:"primaryKey" json:"channel_id"`
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	JoinedAt  time.Time `json:"joined_at"`
}
