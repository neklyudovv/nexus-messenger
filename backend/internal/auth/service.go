package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"nexus-messenger/backend/config"
	"nexus-messenger/backend/internal/user"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	redis *redis.Client
	cfg   *config.Config
}

func NewService(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *Service {
	return &Service{db: db, redis: rdb, cfg: cfg}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (s *Service) Register(username, email, password string) (*user.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Username: username,
		Email:    email,
		Password: string(hash),
	}

	if err := s.db.Create(u).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("user already exists")
		}
		return nil, err
	}
	return u, nil
}

func (s *Service) Login(email, password string) (*user.User, *TokenPair, error) {
	var u user.User
	if err := s.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	pair, err := s.newTokenPair(u.ID)
	if err != nil {
		return nil, nil, err
	}
	return &u, pair, nil
}

func (s *Service) Refresh(refreshToken string) (*TokenPair, error) {
	claims, err := ParseToken(refreshToken, s.cfg.JWTSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if s.redis.Exists(context.Background(), refreshKey(refreshToken)).Val() == 0 {
		return nil, errors.New("refresh token revoked")
	}
	s.redis.Del(context.Background(), refreshKey(refreshToken))

	return s.newTokenPair(claims.UserID)
}

func (s *Service) Logout(refreshToken string) error {
	return s.redis.Del(context.Background(), refreshKey(refreshToken)).Err()
}

func (s *Service) newTokenPair(userID uint) (*TokenPair, error) {
	access, err := generateToken(userID, s.cfg.JWTAccessTTL, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	refresh, err := generateToken(userID, s.cfg.JWTRefreshTTL, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	s.redis.Set(context.Background(), refreshKey(refresh), userID, s.cfg.JWTRefreshTTL+time.Minute)

	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

// refreshKey hashes the token so that long JWT strings are not used as Redis keys.
func refreshKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return "refresh:" + hex.EncodeToString(sum[:])
}
