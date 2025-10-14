package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/model"
	"github.com/protocyber/kelasgo-api/internal/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
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
	// Parse tenant ID if provided
	var tenantID uuid.UUID
	var err error
	if req.TenantID != "" {
		tenantID, err = uuid.Parse(req.TenantID)
		if err != nil {
			log.Error().
				Err(err).
				Str("tenant_id", req.TenantID).
				Str("username", req.Username).
				Msg("Failed to parse tenant ID during login")
			return nil, errors.New("invalid tenant ID format")
		}
	}

	// Get user by username and tenant
	user, err := s.userRepo.GetByUsernameAndTenant(req.Username, tenantID)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("User not found during login attempt")
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		log.Warn().
			Str("user_id", user.ID.String()).
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("Login attempt for deactivated user")
		return nil, errors.New("user account is deactivated")
	}

	// Check password
	if !util.CheckPassword(req.Password, user.PasswordHash) {
		log.Warn().
			Str("user_id", user.ID.String()).
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("Invalid password during login attempt")
		return nil, errors.New("invalid username or password")
	}

	// Get tenant user first
	tenantUser, err := s.tenantUserRepo.GetByTenantAndUser(tenantID, user.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.ID.String()).
			Str("tenant_id", tenantID.String()).
			Msg("Tenant user not found during login")
		return nil, errors.New("user not authorized for this tenant")
	}

	// Get role name from TenantUserRoles
	roleName := ""
	tenantUserRoles, err := s.tenantUserRoleRepo.GetRolesByTenantUser(tenantUser.ID)
	if err == nil && len(tenantUserRoles) > 0 && tenantUserRoles[0].Role != nil {
		roleName = tenantUserRoles[0].Role.Name
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
		log.Error().
			Err(err).
			Str("user_id", user.ID.String()).
			Str("username", user.Username).
			Str("tenant_id", tenantID.String()).
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
		log.Warn().
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("Registration attempt with existing username")
		return nil, errors.New("username already exists")
	}

	// Check if email already exists in this tenant (if provided)
	if req.Email != "" {
		existingUser, _ = s.userRepo.GetByEmailAndTenant(req.Email, tenantID)
		if existingUser != nil {
			log.Warn().
				Str("email", req.Email).
				Str("tenant_id", tenantID.String()).
				Msg("Registration attempt with existing email")
			return nil, errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to hash password during registration")
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
		log.Error().
			Err(err).
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to create user during registration")
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
		log.Error().
			Err(err).
			Str("user_id", user.ID.String()).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to create tenant-user relationship during registration")
		// If tenant-user creation fails, we should delete the user to maintain consistency
		s.userRepo.Delete(user.ID)
		return nil, errors.New("failed to create tenant-user relationship")
	}

	// Create tenant user-role relationship if role is provided
	if req.RoleID != nil {
		tenantUserRole := &model.TenantUserRole{
			TenantUserID: tenantUser.ID,
			RoleID:       *req.RoleID,
		}

		err = s.tenantUserRoleRepo.Create(tenantUserRole)
		if err != nil {
			log.Error().
				Err(err).
				Str("tenant_user_id", tenantUser.ID.String()).
				Str("role_id", req.RoleID.String()).
				Str("tenant_id", tenantID.String()).
				Msg("Failed to create tenant user-role relationship during registration")
			// If tenant user-role creation fails, cleanup user and tenant-user
			s.tenantUserRepo.Delete(tenantUser.ID)
			s.userRepo.Delete(user.ID)
			return nil, errors.New("failed to create tenant user-role relationship")
		}
	}

	return user, nil
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
		TenantID: claims.TenantID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}

	return tokenClaims, nil
}
