package workspace

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(name, description string, ownerID uint) (*Workspace, error) {
	code, err := generateInviteCode()
	if err != nil {
		return nil, err
	}
	ws := &Workspace{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		InviteCode:  code,
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ws).Error; err != nil {
			return err
		}
		return tx.Create(&WorkspaceMember{
			WorkspaceID: ws.ID,
			UserID:      ownerID,
			Role:        RoleAdmin,
			JoinedAt:    time.Now(),
		}).Error
	})
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *Service) GetMine(userID uint) ([]Workspace, error) {
	var workspaces []Workspace
	err := s.db.
		Joins("JOIN workspace_members ON workspace_members.workspace_id = workspaces.id").
		Where("workspace_members.user_id = ?", userID).
		Find(&workspaces).Error
	return workspaces, err
}

func (s *Service) GetByID(id uint) (*Workspace, error) {
	var ws Workspace
	if err := s.db.First(&ws, id).Error; err != nil {
		return nil, errors.New("workspace not found")
	}
	return &ws, nil
}

func (s *Service) Join(inviteCode string, userID uint) (*Workspace, error) {
	var ws Workspace
	if err := s.db.Where("invite_code = ?", inviteCode).First(&ws).Error; err != nil {
		return nil, errors.New("invalid invite code")
	}
	member := &WorkspaceMember{
		WorkspaceID: ws.ID,
		UserID:      userID,
		Role:        RoleMember,
		JoinedAt:    time.Now(),
	}
	if err := s.db.FirstOrCreate(member, WorkspaceMember{WorkspaceID: ws.ID, UserID: userID}).Error; err != nil {
		return nil, err
	}
	return &ws, nil
}

func (s *Service) GetMembers(workspaceID uint) ([]WorkspaceMember, error) {
	var members []WorkspaceMember
	return members, s.db.Where("workspace_id = ?", workspaceID).Find(&members).Error
}

func (s *Service) IsMember(workspaceID, userID uint) bool {
	var count int64
	s.db.Model(&WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count)
	return count > 0
}

func (s *Service) RegenerateInvite(workspaceID, requesterID uint) (string, error) {
	var member WorkspaceMember
	if err := s.db.Where("workspace_id = ? AND user_id = ?", workspaceID, requesterID).First(&member).Error; err != nil {
		return "", errors.New("not a workspace member")
	}
	if member.Role != RoleAdmin {
		return "", errors.New("only admins can regenerate invite")
	}
	var ws Workspace
	if err := s.db.First(&ws, workspaceID).Error; err != nil {
		return "", errors.New("workspace not found")
	}
	code, err := generateInviteCode()
	if err != nil {
		return "", err
	}
	if err := s.db.Model(&ws).Update("invite_code", code).Error; err != nil {
		return "", err
	}
	return code, nil
}

func (s *Service) UpdateMemberRole(workspaceID, requesterID, targetID uint, role MemberRole) error {
	if role != RoleAdmin && role != RoleMember {
		return errors.New("invalid role: must be admin or member")
	}
	if targetID == requesterID {
		return errors.New("cannot change your own role")
	}
	var requester WorkspaceMember
	if err := s.db.Where("workspace_id = ? AND user_id = ?", workspaceID, requesterID).First(&requester).Error; err != nil {
		return errors.New("not a workspace member")
	}
	if requester.Role != RoleAdmin {
		return errors.New("only admins can change roles")
	}
	if err := s.db.Where("workspace_id = ? AND user_id = ?", workspaceID, targetID).First(&WorkspaceMember{}).Error; err != nil {
		return errors.New("target user not in workspace")
	}
	return s.db.Model(&WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, targetID).
		Update("role", role).Error
}

func generateInviteCode() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
