package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"webtest/models"
)

type AuthService interface {
	Register(input models.AuthInput) (*models.User, error)
	Login(input models.AuthInput) (string, time.Time, *models.User, error)
}

type authService struct {
	db          *gorm.DB
	jwtSecret   []byte
	jwtLifetime time.Duration
	bcryptCost  int
}

func NewAuthService(db *gorm.DB, jwtSecret []byte, jwtLifetime time.Duration, bcryptCost int) AuthService {
	return &authService{
		db:          db,
		jwtSecret:   jwtSecret,
		jwtLifetime: jwtLifetime,
		bcryptCost:  bcryptCost,
	}
}

func (s *authService) Register(input models.AuthInput) (*models.User, error) {
	if err := isValidPassword(input.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.bcryptCost)
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("user with this username already exists")
	}

	return &user, nil
}

func (s *authService) Login(input models.AuthInput) (string, time.Time, *models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return "", time.Time{}, nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", time.Time{}, nil, errors.New("invalid username or password")
	}

	expirationTime := time.Now().Add(s.jwtLifetime)
	return "", expirationTime, &user, nil
}
