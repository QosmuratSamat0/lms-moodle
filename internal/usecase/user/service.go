package user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/user"
	"github.com/google/uuid"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(input *user.CreateUserInput) (*user.User, error) {
	hash := sha256.Sum256([]byte(input.Password))
	u := &user.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		Password:  hex.EncodeToString(hash[:]),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      user.RoleUser,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) GetByID(id string) (*user.User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetByEmail(email string) (*user.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *Service) Update(id string, input *user.UpdateUserInput) (*user.User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.FirstName != nil {
		u.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		u.LastName = *input.LastName
	}
	if input.Active != nil {
		u.Active = *input.Active
	}
	u.UpdatedAt = time.Now()
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) List(skip, take int) ([]*user.User, error) {
	return s.repo.List(skip, take)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) Login(email, password string) (string, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	hash := sha256.Sum256([]byte(password))
	hashedPassword := hex.EncodeToString(hash[:])

	if hashedPassword != u.Password {
		return "", errors.New("invalid password")
	}

	if !u.Active {
		return "", errors.New("user is inactive")
	}

	token := email + ":" + strconv.FormatInt(time.Now().Unix(), 10)
	return token, nil
}

// GetByEmailAndPassword валидирует пользователя по email и паролю
func (s *Service) GetByEmailAndPassword(email, password string) (*user.User, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Хешируем пароль для сравнения
	hash := sha256.Sum256([]byte(password))
	hashedPassword := hex.EncodeToString(hash[:])

	if hashedPassword != u.Password {
		return nil, errors.New("invalid password")
	}

	if !u.Active {
		return nil, errors.New("user is inactive")
	}

	return u, nil
}
