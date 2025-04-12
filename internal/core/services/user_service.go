package services

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/pkg/apperror"
	"socket/pkg/util"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo ports.UserRepository
	jwt  util.JWTUtils
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(userName, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.InternalServerError(err, "cannot use this password")
	}
	user := domain.User{Username: userName, Password: string(hashedPassword)}
	err = s.repo.Create(&user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(userName, password string) (*domain.User, string, error) {
	user, err := s.repo.GetUserByUsername(userName)
	if err != nil {
		return nil, "", err
	}
	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", apperror.InternalServerError(err, "cannot use this password")
	}
	if string(hashedPasswordByte) != user.Password {
		return nil, "", apperror.BadRequestError(err, "username or password is incorrect")
	}
	jwt, err := s.jwt.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, jwt, nil
}
