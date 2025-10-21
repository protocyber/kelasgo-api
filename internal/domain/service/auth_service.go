package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/dto"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/domain/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// AuthService interface defines authentication service methods
type AuthService interface {
	Login(c context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(c context.Context, req dto.RegisterRequest) (*model.User, error)
	SelectTenant(c context.Context, userID uuid.UUID, req dto.TenantSelectionRequest) (*dto.TenantSelectionResponse, error)
	GetUserTenants(c context.Context, userID uuid.UUID) ([]model.TenantUser, error)
	ChangePassword(c context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error
	ValidateToken(c context.Context, token string) (*dto.TokenClaims, error)
}

// authService implements AuthService
type authService struct {
	userRepo           repository.UserRepository
	roleRepo           repository.RoleRepository
	tenantUserRepo     repository.TenantUserRepository
	tenantUserRoleRepo repository.TenantUserRoleRepository
	jwtService         *util.JWTService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tenantUserRepo repository.TenantUserRepository,
	tenantUserRoleRepo repository.TenantUserRoleRepository,
	jwtService *util.JWTService,
) AuthService {
	return &authService{
		userRepo:           userRepo,
		roleRepo:           roleRepo,
		tenantUserRepo:     tenantUserRepo,
		tenantUserRoleRepo: tenantUserRoleRepo,
		jwtService:         jwtService,
	}
}

func (s *authService) Login(c context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Get user by email globally (no tenant context needed)
	user, err := s.userRepo.GetByEmailGlobal(c, req.Email)
	if err != nil {
		logger.Error().
			Err(err).
			Str("email", req.Email).
			Msg("User not found during login attempt")
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		logger.Warn().
			Str("user_id", user.ID.String()).
			Str("email", req.Email).
			Msg("Login attempt for deactivated user")
		return nil, errors.New("user account is deactivated")
	}

	// Check password
	if !util.CheckPassword(req.Password, user.PasswordHash) {
		logger.Warn().
			Str("user_id", user.ID.String()).
			Str("email", req.Email).
			Msg("Invalid password during login attempt")
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token without tenant context (user can select tenant later)
	token, expiresAt, err := s.jwtService.GenerateToken(
		user.ID,
		uuid.Nil, // No tenant selected yet
		user.Username,
		user.Email,
		"", // No role yet
	)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", user.ID.String()).
			Str("email", req.Email).
			Msg("Failed to generate JWT token during login")
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
			TenantID: nil, // No tenant selected yet
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     "", // No role yet
		},
	}, nil
}

func (s *authService) Register(c context.Context, req dto.RegisterRequest) (*model.User, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Check if username already exists globally
	existingUser, _ := s.userRepo.GetByUsername(c, req.Username)
	if existingUser != nil {
		logger.Warn().
			Str("username", req.Username).
			Msg("Registration attempt with existing username")
		return nil, errors.New("username already exists")
	}

	// Check if email already exists globally
	existingUser, _ = s.userRepo.GetByEmailGlobal(c, req.Email)
	if existingUser != nil {
		logger.Warn().
			Str("email", req.Email).
			Msg("Registration attempt with existing email")
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		logger.Error().
			Err(err).
			Str("email", req.Email).
			Msg("Failed to hash password during registration")
		return nil, errors.New("failed to hash password")
	}

	// Create user without tenant context
	user := &model.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
		FullName:     req.FullName,
		Phone:        &req.Phone,
		IsActive:     true,
	}

	err = s.userRepo.Create(c, user)
	if err != nil {
		logger.Error().
			Err(err).
			Str("email", req.Email).
			Str("username", req.Username).
			Msg("Failed to create user during registration")
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *authService) SelectTenant(c context.Context, userID uuid.UUID, req dto.TenantSelectionRequest) (*dto.TenantSelectionResponse, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Parse tenant ID
	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("tenant_id", req.TenantID).
			Str("user_id", userID.String()).
			Msg("Failed to parse tenant ID during tenant selection")
		return nil, errors.New("invalid tenant ID format")
	}

	// Get user
	user, err := s.userRepo.GetByID(c, userID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("User not found during tenant selection")
		return nil, errors.New("user not found")
	}

	// Check if user belongs to this tenant
	tenantUser, err := s.tenantUserRepo.GetByTenantAndUser(c, tenantID, userID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Str("tenant_id", tenantID.String()).
			Msg("User not authorized for this tenant")
		return nil, errors.New("user not authorized for this tenant")
	}

	// Get role name from TenantUserRoles
	roleName := ""
	tenantUserRoles, err := s.tenantUserRoleRepo.GetRolesByTenantUser(c, tenantUser.ID)
	if err == nil && len(tenantUserRoles) > 0 && tenantUserRoles[0].Role != nil {
		roleName = tenantUserRoles[0].Role.Name
	}

	// Generate JWT token with tenant context
	token, expiresAt, err := s.jwtService.GenerateToken(
		user.ID,
		tenantID,
		user.Username,
		user.Email,
		roleName,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", user.ID.String()).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to generate JWT token during tenant selection")
		return nil, errors.New("failed to generate token")
	}

	// TODO: Implement refresh token logic
	refreshToken := token // For now, use same token

	return &dto.TenantSelectionResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: dto.UserInfo{
			ID:       user.ID,
			TenantID: &tenantID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     roleName,
		},
	}, nil
}

func (s *authService) GetUserTenants(c context.Context, userID uuid.UUID) ([]model.TenantUser, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	tenantUsers, err := s.userRepo.GetUserTenants(c, userID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to get user tenants")
		return nil, errors.New("failed to get user tenants")
	}
	return tenantUsers, nil
}

func (s *authService) ChangePassword(c context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Get user
	user, err := s.userRepo.GetByID(c, userID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("User not found during password change")
		return errors.New("user not found")
	}

	// Check current password
	if !util.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		logger.Warn().
			Str("user_id", userID.String()).
			Str("username", user.Username).
			Msg("Incorrect current password during password change")
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to hash new password during password change")
		return errors.New("failed to hash new password")
	}

	// Update password
	user.PasswordHash = hashedPassword

	err = s.userRepo.Update(c, user)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to update password in database")
		return errors.New("failed to update password")
	}

	return nil
}

func (s *authService) ValidateToken(c context.Context, token string) (*dto.TokenClaims, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Validate JWT token using JWT service
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Token validation failed")
		return nil, errors.New("invalid token")
	}

	// Convert JWT claims to DTO claims
	tokenClaims := &dto.TokenClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}

	// Handle optional tenant ID
	if claims.TenantID != uuid.Nil {
		tokenClaims.TenantID = &claims.TenantID
	}

	return tokenClaims, nil
}
