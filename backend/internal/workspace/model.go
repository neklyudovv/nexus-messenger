package workspace

import "time"

type Workspace struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	OwnerID     uint      `gorm:"not null" json:"owner_id"`
	InviteCode  string    `gorm:"uniqueIndex;not null" json:"invite_code"`
	CreatedAt   time.Time `json:"created_at"`
}

type MemberRole string

const (
	RoleAdmin  MemberRole = "admin"
	RoleMember MemberRole = "member"
)

type WorkspaceMember struct {
	WorkspaceID uint       `gorm:"primaryKey" json:"workspace_id"`
	UserID      uint       `gorm:"primaryKey" json:"user_id"`
	Role        MemberRole `gorm:"default:member" json:"role"`
	JoinedAt    time.Time  `json:"joined_at"`
}
