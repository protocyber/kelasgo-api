package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/dto"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/domain/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
)

// AuthService interface defines authentication service methods
type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req dto.RegisterRequest) (*model.User, error)
	SelectTenant(userID uuid.UUID, req dto.TenantSelectionRequest) (*dto.TenantSelectionResponse, error)
	GetUserTenants(userID uuid.UUID) ([]model.TenantUser, error)
	ChangePassword(userID uuid.UUID, req dto.ChangePasswordRequest) error
	ValidateToken(token string) (*dto.TokenClaims, error)
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

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email globally (no tenant context needed)
	user, err := s.userRepo.GetByEmailGlobal(req.Email)
	if err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("User not found during login attempt")
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		log.Warn().
			Str("user_id", user.ID.String()).
			Str("email", req.Email).
			Msg("Login attempt for deactivated user")
		return nil, errors.New("user account is deactivated")
	}

	// Check password
	if !util.CheckPassword(req.Password, user.PasswordHash) {
		log.Warn().
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
		log.Error().
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

func (s *authService) Register(req dto.RegisterRequest) (*model.User, error) {
	// Check if username already exists globally
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		log.Warn().
			Str("username", req.Username).
			Msg("Registration attempt with existing username")
		return nil, errors.New("username already exists")
	}

	// Check if email already exists globally
	existingUser, _ = s.userRepo.GetByEmailGlobal(req.Email)
	if existingUser != nil {
		log.Warn().
			Str("email", req.Email).
			Msg("Registration attempt with existing email")
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Error().
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

	err = s.userRepo.Create(user)
	if err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Str("username", req.Username).
			Msg("Failed to create user during registration")
		return nil, errors.New("failed to create user")
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("email", req.Email).
		Str("username", req.Username).
		Msg("User registered successfully without tenant")

	return user, nil
}

func (s *authService) SelectTenant(userID uuid.UUID, req dto.TenantSelectionRequest) (*dto.TenantSelectionResponse, error) {
	// Parse tenant ID
	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", req.TenantID).
			Str("user_id", userID.String()).
			Msg("Failed to parse tenant ID during tenant selection")
		return nil, errors.New("invalid tenant ID format")
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("User not found during tenant selection")
		return nil, errors.New("user not found")
	}

	// Check if user belongs to this tenant
	tenantUser, err := s.tenantUserRepo.GetByTenantAndUser(tenantID, userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Str("tenant_id", tenantID.String()).
			Msg("User not authorized for this tenant")
		return nil, errors.New("user not authorized for this tenant")
	}

	// Get role name from TenantUserRoles
	roleName := ""
	tenantUserRoles, err := s.tenantUserRoleRepo.GetRolesByTenantUser(tenantUser.ID)
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
		log.Error().
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

func (s *authService) GetUserTenants(userID uuid.UUID) ([]model.TenantUser, error) {
	tenantUsers, err := s.userRepo.GetUserTenants(userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to get user tenants")
		return nil, errors.New("failed to get user tenants")
	}
	return tenantUsers, nil
}

func (s *authService) ChangePassword(userID uuid.UUID, req dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("User not found during password change")
		return errors.New("user not found")
	}

	// Check current password
	if !util.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		log.Warn().
			Str("user_id", userID.String()).
			Str("username", user.Username).
			Msg("Incorrect current password during password change")
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to hash new password during password change")
		return errors.New("failed to hash new password")
	}

	// Update password
	user.PasswordHash = hashedPassword

	err = s.userRepo.Update(user)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to update password in database")
		return errors.New("failed to update password")
	}

	log.Info().
		Str("user_id", userID.String()).
		Str("username", user.Username).
		Msg("Password changed successfully")

	return nil
}

func (s *authService) ValidateToken(token string) (*dto.TokenClaims, error) {
	// Validate JWT token using JWT service
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		log.Error().
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
