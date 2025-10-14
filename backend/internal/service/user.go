package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/auth"
	"github.com/v-ui/backend/internal/config"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

type UserService struct {
	db  *db.Database
	cfg *config.Config
	log *zap.Logger
}

func NewUserService(database *db.Database, cfg *config.Config, log *zap.Logger) *UserService {
	return &UserService{
		db:  database,
		cfg: cfg,
		log: log,
	}
}

func (s *UserService) Create(username, email, password string, role models.Role) (*models.User, error) {
	passwordHash, err := auth.HashPassword(password, s.cfg.Auth.BcryptCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		Enabled:      true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) List(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) Update(id uuid.UUID, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.User{}, id).Error
}

func (s *UserService) ChangePassword(id uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if !auth.CheckPassword(oldPassword, user.PasswordHash) {
		return errors.New("invalid old password")
	}

	passwordHash, err := auth.HashPassword(newPassword, s.cfg.Auth.BcryptCost)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	return s.db.Save(user).Error
}

