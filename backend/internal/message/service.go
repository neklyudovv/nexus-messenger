package message

import (
	"errors"

	"gorm.io/gorm"
)

const pageSize = 50

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GetHistory returns up to 50 messages for a channel in chronological order.
// Pass beforeID to paginate backward (returns messages with id < beforeID).
func (s *Service) GetHistory(channelID, beforeID uint) ([]Message, error) {
	var msgs []Message
	q := s.db.Where("channel_id = ? AND is_deleted = false", channelID)
	if beforeID > 0 {
		q = q.Where("id < ?", beforeID)
	}
	err := q.Order("id DESC").Limit(pageSize).Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}

func (s *Service) Create(channelID, userID uint, content string) (*Message, error) {
	msg := &Message{
		ChannelID: channelID,
		UserID:    userID,
		Content:   content,
		Type:      TypeText,
	}
	return msg, s.db.Create(msg).Error
}

// Delete soft-deletes a message. The requester must be the author or a workspace admin.
// Returns the channel ID on success so callers can broadcast the event.
func (s *Service) Delete(msgID, userID uint) (channelID uint, err error) {
	var msg Message
	if err := s.db.First(&msg, msgID).Error; err != nil {
		return 0, ErrNotFound
	}
	if msg.UserID != userID {
		var count int64
		s.db.Raw(
			`SELECT COUNT(*) FROM workspace_members wm
			 INNER JOIN channels c ON c.workspace_id = wm.workspace_id
			 WHERE c.id = ? AND wm.user_id = ? AND wm.role = 'admin'`,
			msg.ChannelID, userID,
		).Scan(&count)
		if count == 0 {
			return 0, ErrForbidden
		}
	}
	msg.IsDeleted = true
	msg.Content = ""
	return msg.ChannelID, s.db.Save(&msg).Error
}
