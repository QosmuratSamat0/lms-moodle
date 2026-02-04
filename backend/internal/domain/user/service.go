package user

import (
	"context"
	"errors"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the user service interface
type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, role string, req *UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID uuid.UUID, req *ChangePasswordRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*UserResponse, error)
	List(ctx context.Context, page, limit int) (*UserListResponse, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
	Activate(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo       Repository
	jwtManager *utils.JWTManager
}

// NewService creates a new user service
func NewService(repo Repository, jwtManager *utils.JWTManager) Service {
	return &service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	// Check if email already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.Wrap(err, "check existing email")
	}
	if existing != nil {
		return nil, errorx.NewConflictError("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errorx.Wrap(err, "hash password")
	}

	// Create user
	user := &User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		IsActive:     true,
	}

	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "begin transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Create user record within transaction
	if err = s.repo.CreateWithTx(ctx, tx, user); err != nil {
		return nil, errorx.Wrap(err, "create user")
	}

	// Create role-specific profile
	if req.Role != "admin" {
		if err = s.repo.CreateProfile(ctx, tx, user.ID, req.Role, req); err != nil {
			return nil, errorx.Wrap(err, "create profile")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errorx.Wrap(err, "commit transaction")
	}

	return s.toUserResponse(user, req.FirstName, req.LastName, req.GroupName, req.Department), nil
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewUnauthorizedError("invalid email or password")
		}
		return nil, errorx.Wrap(err, "get user by email")
	}

	if !user.IsActive {
		return nil, errorx.NewForbiddenError("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errorx.NewUnauthorizedError("invalid email or password")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errorx.Wrap(err, "generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "generate refresh token")
	}

	// Get full profile
	profile, _ := s.repo.GetProfile(ctx, user.ID, user.Role)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *s.profileToUserResponse(profile),
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*LoginResponse, error) {
	userID, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errorx.NewUnauthorizedError("invalid refresh token")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("user")
		}
		return nil, errorx.Wrap(err, "get user")
	}

	if !user.IsActive {
		return nil, errorx.NewForbiddenError("account is deactivated")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errorx.Wrap(err, "generate access token")
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "generate refresh token")
	}

	profile, _ := s.repo.GetProfile(ctx, user.ID, user.Role)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         *s.profileToUserResponse(profile),
	}, nil
}

func (s *service) GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("user")
		}
		return nil, errorx.Wrap(err, "get user")
	}

	profile, err := s.repo.GetProfile(ctx, userID, user.Role)
	if err != nil {
		return nil, errorx.Wrap(err, "get profile")
	}

	return s.profileToUserResponse(profile), nil
}

func (s *service) UpdateProfile(ctx context.Context, userID uuid.UUID, role string, req *UpdateProfileRequest) error {
	if err := s.repo.UpdateProfile(ctx, userID, role, req); err != nil {
		return errorx.Wrap(err, "update profile")
	}
	return nil
}

func (s *service) ChangePassword(ctx context.Context, userID uuid.UUID, req *ChangePasswordRequest) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("user")
		}
		return errorx.Wrap(err, "get user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errorx.NewBadRequestError("incorrect current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errorx.Wrap(err, "hash password")
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.repo.Update(ctx, user); err != nil {
		return errorx.Wrap(err, "update user")
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("user")
		}
		return nil, errorx.Wrap(err, "get user")
	}

	profile, _ := s.repo.GetProfile(ctx, id, user.Role)
	return s.profileToUserResponse(profile), nil
}

func (s *service) List(ctx context.Context, page, limit int) (*UserListResponse, error) {
	offset := (page - 1) * limit
	users, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list users")
	}

	var userResponses []UserResponse
	for _, u := range users {
		profile, _ := s.repo.GetProfile(ctx, u.ID, u.Role)
		userResponses = append(userResponses, *s.profileToUserResponse(profile))
	}

	return &UserListResponse{
		Users: userResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) Deactivate(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("user")
		}
		return errorx.Wrap(err, "get user")
	}

	user.IsActive = false
	if err := s.repo.Update(ctx, user); err != nil {
		return errorx.Wrap(err, "deactivate user")
	}

	return nil
}

func (s *service) Activate(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("user")
		}
		return errorx.Wrap(err, "get user")
	}

	user.IsActive = true
	if err := s.repo.Update(ctx, user); err != nil {
		return errorx.Wrap(err, "activate user")
	}

	return nil
}

// Helper methods

func (s *service) toUserResponse(user *User, firstName, lastName, groupName, department *string) *UserResponse {
	return &UserResponse{
		ID:         user.ID.String(),
		Email:      user.Email,
		Role:       user.Role,
		IsActive:   user.IsActive,
		FirstName:  firstName,
		LastName:   lastName,
		GroupName:  groupName,
		Department: department,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *service) profileToUserResponse(profile *UserProfile) *UserResponse {
	if profile == nil {
		return nil
	}

	resp := &UserResponse{
		ID:        profile.User.ID.String(),
		Email:     profile.User.Email,
		Role:      profile.User.Role,
		IsActive:  profile.User.IsActive,
		CreatedAt: profile.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	switch {
	case profile.Student != nil:
		resp.FirstName = profile.Student.FirstName
		resp.LastName = profile.Student.LastName
		resp.GroupName = profile.Student.GroupName
	case profile.Teacher != nil:
		resp.FirstName = profile.Teacher.FirstName
		resp.LastName = profile.Teacher.LastName
		resp.Department = profile.Teacher.Department
	case profile.Manager != nil:
		resp.FirstName = profile.Manager.FirstName
		resp.LastName = profile.Manager.LastName
	}

	return resp
}

// Ensure the service implements Service interface
var _ Service = (*service)(nil)
