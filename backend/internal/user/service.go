package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	onlineTTL = 35 * time.Second
	typingTTL = 3 * time.Second
)

type Service struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{db: db, redis: rdb}
}

func (s *Service) GetAll() ([]User, error) {
	var users []User
	return users, s.db.Find(&users).Error
}

func (s *Service) GetByID(id uint) (*User, error) {
	var u User
	if err := s.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) Update(id uint, username, avatarURL string) (*User, error) {
	var u User
	if err := s.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	if username != "" {
		u.Username = username
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
	return &u, s.db.Save(&u).Error
}

func (s *Service) SetOnline(userID uint) {
	if err := s.redis.Set(context.Background(), onlineKey(userID), 1, onlineTTL).Err(); err != nil {
		log.Printf("redis SetOnline user=%d: %v", userID, err)
	}
}

func (s *Service) SetOffline(userID uint) {
	if err := s.redis.Del(context.Background(), onlineKey(userID)).Err(); err != nil {
		log.Printf("redis SetOffline user=%d: %v", userID, err)
	}
}

func (s *Service) IsOnline(userID uint) bool {
	return s.redis.Exists(context.Background(), onlineKey(userID)).Val() == 1
}

func (s *Service) SetTyping(channelID, userID uint) {
	key := fmt.Sprintf("channel:%d:typing:%d", channelID, userID)
	if err := s.redis.Set(context.Background(), key, 1, typingTTL).Err(); err != nil {
		log.Printf("redis SetTyping ch=%d user=%d: %v", channelID, userID, err)
	}
}

func onlineKey(userID uint) string {
	return fmt.Sprintf("user:%d:online", userID)
}
