package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/model"
	"github.com/protocyber/kelasgo-api/internal/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// AuthService interface defines authentication service methods
type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(tenantID uuid.UUID, req dto.CreateUserRequest) (*model.User, error)
	ChangePassword(userID uuid.UUID, req dto.ChangePasswordRequest) error
	ValidateToken(token string) (*dto.TokenClaims, error)
}

// authService implements AuthService
type authService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	tenantUserRepo repository.TenantUserRepository
	userRoleRepo   repository.UserRoleRepository
	jwtService     *util.JWTService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tenantUserRepo repository.TenantUserRepository,
	userRoleRepo repository.UserRoleRepository,
	jwtService *util.JWTService,
) AuthService {
	return &authService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		tenantUserRepo: tenantUserRepo,
		userRoleRepo:   userRoleRepo,
		jwtService:     jwtService,
	}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Parse tenant ID if provided
	var tenantID uuid.UUID
	var err error
	if req.TenantID != "" {
		tenantID, err = uuid.Parse(req.TenantID)
		if err != nil {
			return nil, errors.New("invalid tenant ID format")
		}
	}

	// Get user by username and tenant
	user, err := s.userRepo.GetByUsernameAndTenant(req.Username, tenantID)
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

	// Get role name from UserRoles
	roleName := ""
	if len(user.UserRoles) > 0 && user.UserRoles[0].Role != nil {
		roleName = user.UserRoles[0].Role.Name
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtService.GenerateToken(
		user.ID,
		tenantID,
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
			TenantID: tenantID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     roleName,
		},
	}, nil
}

func (s *authService) Register(tenantID uuid.UUID, req dto.CreateUserRequest) (*model.User, error) {
	// Check if username already exists in this tenant
	existingUser, _ := s.userRepo.GetByUsernameAndTenant(req.Username, tenantID)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists in this tenant (if provided)
	if req.Email != "" {
		existingUser, _ = s.userRepo.GetByEmailAndTenant(req.Email, tenantID)
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
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
		FullName:     req.FullName,
		Birthplace:   req.Birthplace,
		Birthday:     req.Birthday,
		Gender:       (*model.Gender)(req.Gender),
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

	// Create tenant-user relationship
	tenantUser := &model.TenantUser{
		TenantID: tenantID,
		UserID:   user.ID,
		IsActive: user.IsActive,
	}

	err = s.tenantUserRepo.Create(tenantUser)
	if err != nil {
		// If tenant-user creation fails, we should delete the user to maintain consistency
		s.userRepo.Delete(user.ID)
		return nil, errors.New("failed to create tenant-user relationship")
	}

	// Create user-role relationship if role is provided
	if req.RoleID != nil {
		userRole := &model.UserRole{
			UserID: user.ID,
			RoleID: *req.RoleID,
		}

		err = s.userRoleRepo.Create(userRole)
		if err != nil {
			// If user-role creation fails, cleanup user and tenant-user
			s.tenantUserRepo.Delete(tenantUser.ID)
			s.userRepo.Delete(user.ID)
			return nil, errors.New("failed to create user-role relationship")
		}
	}

	return user, nil
}

func (s *authService) ChangePassword(userID uuid.UUID, req dto.ChangePasswordRequest) error {
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

	err = s.userRepo.Update(user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *authService) ValidateToken(token string) (*dto.TokenClaims, error) {
	// Validate JWT token using JWT service
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Convert JWT claims to DTO claims
	tokenClaims := &dto.TokenClaims{
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}

	return tokenClaims, nil
}
