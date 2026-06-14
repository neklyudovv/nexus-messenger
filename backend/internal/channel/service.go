package channel

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

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

func (s *Service) GetForWorkspace(workspaceID uint) ([]Channel, error) {
	var channels []Channel
	err := s.db.
		Where("workspace_id = ? AND type != ?", workspaceID, TypeDM).
		Find(&channels).Error
	return channels, err
}

func (s *Service) GetAllDMs(userID uint) ([]Channel, error) {
	var channels []Channel
	err := s.db.
		Joins("JOIN channel_members ON channel_members.channel_id = channels.id AND channel_members.user_id = ?", userID).
		Where("channels.type = ?", TypeDM).
		Find(&channels).Error
	return channels, err
}

func (s *Service) GetByID(id uint) (*Channel, error) {
	var ch Channel
	if err := s.db.First(&ch, id).Error; err != nil {
		return nil, ErrNotFound
	}
	return &ch, nil
}

func (s *Service) GetWorkspaceMemberIDs(workspaceID uint) ([]uint, error) {
	var ids []uint
	err := s.db.Table("workspace_members").
		Where("workspace_id = ?", workspaceID).
		Pluck("user_id", &ids).Error
	return ids, err
}

func (s *Service) Create(workspaceID uint, name, description string, chType ChannelType, createdBy uint) (*Channel, error) {
	ch := &Channel{
		WorkspaceID: workspaceID,
		Name:        name,
		Description: description,
		Type:        chType,
		CreatedBy:   createdBy,
	}
	if err := s.db.Create(ch).Error; err != nil {
		return nil, err
	}
	return ch, nil
}

// Delete removes the channel. Requester must be workspace admin or channel creator.
func (s *Service) Delete(id, userID uint) error {
	var ch Channel
	if err := s.db.First(&ch, id).Error; err != nil {
		return ErrNotFound
	}
	var member struct{ Role string }
	s.db.Table("workspace_members").
		Select("role").
		Where("workspace_id = ? AND user_id = ?", ch.WorkspaceID, userID).
		Scan(&member)
	if member.Role != "admin" && ch.CreatedBy != userID {
		return ErrForbidden
	}
	return s.db.Delete(&Channel{}, id).Error
}

func (s *Service) Join(channelID, userID uint) error {
	ch, err := s.GetByID(channelID)
	if err != nil {
		return ErrNotFound
	}
	if ch.Type == TypeDM {
		return errors.New("cannot join DM directly")
	}
	member := &ChannelMember{ChannelID: channelID, UserID: userID, JoinedAt: time.Now()}
	return s.db.FirstOrCreate(member, ChannelMember{ChannelID: channelID, UserID: userID}).Error
}

func (s *Service) Leave(channelID, userID uint) error {
	return s.db.Where("channel_id = ? AND user_id = ?", channelID, userID).
		Delete(&ChannelMember{}).Error
}

// AddMember adds a user to a channel. Requester must be channel creator or workspace admin.
func (s *Service) AddMember(channelID, requesterID, targetUserID uint) error {
	ch, err := s.GetByID(channelID)
	if err != nil {
		return ErrNotFound
	}
	if ch.Type == TypeDM {
		return errors.New("cannot add members to DM")
	}
	if ch.CreatedBy != requesterID {
		var count int64
		s.db.Table("workspace_members").
			Where("workspace_id = ? AND user_id = ? AND role = 'admin'", ch.WorkspaceID, requesterID).
			Count(&count)
		if count == 0 {
			return ErrForbidden
		}
	}
	member := &ChannelMember{ChannelID: channelID, UserID: targetUserID, JoinedAt: time.Now()}
	return s.db.FirstOrCreate(member, ChannelMember{ChannelID: channelID, UserID: targetUserID}).Error
}

func (s *Service) GetMembers(channelID uint) ([]ChannelMember, error) {
	var members []ChannelMember
	return members, s.db.Where("channel_id = ?", channelID).Find(&members).Error
}

func (s *Service) IsMember(channelID, userID uint) (bool, error) {
	ch, err := s.GetByID(channelID)
	if err != nil {
		return false, ErrNotFound
	}
	if ch.Type == TypeDM {
		var count int64
		if err := s.db.Model(&ChannelMember{}).
			Where("channel_id = ? AND user_id = ?", channelID, userID).
			Count(&count).Error; err != nil {
			return false, err
		}
		return count > 0, nil
	}
	var count int64
	if err := s.db.Table("workspace_members").
		Where("workspace_id = ? AND user_id = ?", ch.WorkspaceID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// OpenDM finds or creates a direct message channel between two users in the given workspace.
func (s *Service) OpenDM(workspaceID, userID, targetID uint) (*Channel, error) {
	var target struct{ ID uint }
	if err := s.db.Table("users").Select("id").First(&target, targetID).Error; err != nil {
		return nil, errors.New("target user not found")
	}

	var ch Channel
	err := s.db.Transaction(func(tx *gorm.DB) error {
		e := tx.
			Joins("JOIN channel_members m1 ON m1.channel_id = channels.id AND m1.user_id = ?", userID).
			Joins("JOIN channel_members m2 ON m2.channel_id = channels.id AND m2.user_id = ?", targetID).
			Where("channels.type = ? AND channels.workspace_id = ?", TypeDM, workspaceID).
			First(&ch).Error
		if e == nil {
			return nil // already exists
		}
		name := fmt.Sprintf("dm-%d-%d", min(userID, targetID), max(userID, targetID))
		ch = Channel{WorkspaceID: workspaceID, Name: name, Type: TypeDM, CreatedBy: userID}
		if e := tx.Create(&ch).Error; e != nil {
			return e
		}
		if e := tx.Create(&ChannelMember{ChannelID: ch.ID, UserID: userID, JoinedAt: time.Now()}).Error; e != nil {
			return e
		}
		return tx.Create(&ChannelMember{ChannelID: ch.ID, UserID: targetID, JoinedAt: time.Now()}).Error
	})
	return &ch, err
}
