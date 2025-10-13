package service

import (
	"errors"
	"time"

	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/model"
	"github.com/protocyber/kelasgo-api/internal/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// AuthService interface defines authentication service methods
type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req dto.CreateUserRequest) (*model.User, error)
	ChangePassword(userID uint, req dto.ChangePasswordRequest) error
}

// authService implements AuthService
type authService struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	jwtService *util.JWTService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	jwtService *util.JWTService,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Check password
	if !util.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	// Get role name
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtService.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		roleName,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// TODO: Implement refresh token logic
	refreshToken := token // For now, use same token

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     roleName,
		},
	}, nil
}

func (s *authService) Register(req dto.CreateUserRequest) (*model.User, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists (if provided)
	if req.Email != "" {
		existingUser, _ = s.userRepo.GetByEmail(req.Email)
		if existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		RoleID:       req.RoleID,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
		FullName:     req.FullName,
		Gender:       req.Gender,
		DateOfBirth:  req.DateOfBirth,
		Phone:        req.Phone,
		Address:      req.Address,
		IsActive:     true,
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *authService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check current password
	if !util.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Update password
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
