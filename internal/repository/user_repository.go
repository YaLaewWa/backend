package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/pkg/apperror"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {

	exitsUser := r.db.Model(&domain.User{}).Where("username = ?", user.Username).First(&domain.User{})
	if exitsUser.RowsAffected > 0 {
		return apperror.ConflictError(errors.New("username already exists"), "username already exists")
	}

	err := r.db.Create(&user).Error

	if err != nil {
		return apperror.InternalServerError(err, "cannot save user")
	}

	return nil
}

func (r *UserRepository) GetUserByUsername(userName string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("username = ?", userName).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(result.Error, "user not found")
		}
		return nil, apperror.InternalServerError(result.Error, "cannot get user")
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(result.Error, "user not found")
		}
		return nil, apperror.InternalServerError(result.Error, "cannot get user")
	}

	return &user, nil
}
