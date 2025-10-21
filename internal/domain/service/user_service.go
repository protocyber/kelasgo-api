package service

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/dto"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/domain/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// UserService interface defines user service methods
type UserService interface {
	Create(c context.Context, tenantID uuid.UUID, req dto.CreateUserRequest) (*model.User, error)
	GetByID(c context.Context, id uuid.UUID) (*model.User, error)
	Update(c context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*model.User, error)
	Delete(c context.Context, id uuid.UUID) error
	BulkDelete(c context.Context, tenantID uuid.UUID, ids []uuid.UUID) error
	List(c context.Context, tenantID uuid.UUID, params dto.UserQueryParams) ([]model.User, *dto.PaginationMeta, error)
}

// userService implements UserService
type userService struct {
	userRepo           repository.UserRepository
	roleRepo           repository.RoleRepository
	tenantUserRepo     repository.TenantUserRepository
	tenantUserRoleRepo repository.TenantUserRoleRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tenantUserRepo repository.TenantUserRepository,
	tenantUserRoleRepo repository.TenantUserRoleRepository,
) UserService {
	return &userService{
		userRepo:           userRepo,
		roleRepo:           roleRepo,
		tenantUserRepo:     tenantUserRepo,
		tenantUserRoleRepo: tenantUserRoleRepo,
	}
}

func (s *userService) Create(c context.Context, tenantID uuid.UUID, req dto.CreateUserRequest) (*model.User, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Check if username already exists within tenant
	existingUser, _ := s.userRepo.GetByUsernameAndTenant(c, req.Username, tenantID)
	if existingUser != nil {
		logger.Warn().
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("Username already exists within tenant")
		return nil, errors.New("username already exists")
	}

	// Check if email already exists within tenant (if provided)
	if req.Email != "" {
		existingUser, _ = s.userRepo.GetByEmailAndTenant(c, req.Email, tenantID)
		if existingUser != nil {
			logger.Warn().
				Str("email", req.Email).
				Str("tenant_id", tenantID.String()).
				Msg("User creation attempt with existing email")
			return nil, errors.New("email already exists")
		}
	}

	// Validate role if provided
	if req.RoleID != nil {
		_, err := s.roleRepo.GetByID(c, *req.RoleID)
		if err != nil {
			logger.Error().
				Err(err).
				Str("role_id", req.RoleID.String()).
				Str("tenant_id", tenantID.String()).
				Msg("Invalid role ID provided during user creation")
			return nil, errors.New("invalid role ID")
		}
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		logger.Error().
			Err(err).
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to hash password during user creation")
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

	err = s.userRepo.Create(c, user)
	if err != nil {
		logger.Error().
			Err(err).
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to create user in database")
		return nil, errors.New("failed to create user")
	}

	// Create tenant-user relationship
	tenantUser := &model.TenantUser{
		TenantID: tenantID,
		UserID:   user.ID,
		IsActive: user.IsActive,
	}

	err = s.tenantUserRepo.Create(c, tenantUser)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", user.ID.String()).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to create tenant-user relationship")
		// If tenant-user creation fails, we should delete the user to maintain consistency
		s.userRepo.Delete(c, user.ID)
		return nil, errors.New("failed to create tenant-user relationship")
	}

	// Create tenant user-role relationship if role is provided
	if req.RoleID != nil {
		tenantUserRole := &model.TenantUserRole{
			TenantUserID: tenantUser.ID,
			RoleID:       *req.RoleID,
		}

		err = s.tenantUserRoleRepo.Create(c, tenantUserRole)
		if err != nil {
			logger.Error().
				Err(err).
				Str("tenant_user_id", tenantUser.ID.String()).
				Str("role_id", req.RoleID.String()).
				Msg("Failed to create tenant user-role relationship")
			// If tenant user-role creation fails, cleanup user and tenant-user
			s.tenantUserRepo.Delete(c, tenantUser.ID)
			s.userRepo.Delete(c, user.ID)
			return nil, errors.New("failed to create tenant user-role relationship")
		}
	}

	return user, nil
}

func (s *userService) GetByID(c context.Context, id uuid.UUID) (*model.User, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	user, err := s.userRepo.GetByID(c, id)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Failed to get user by ID")
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) Update(c context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*model.User, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Get existing user
	user, err := s.userRepo.GetByID(c, id)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("User not found during update")
		return nil, err
	}

	// Get the user's tenant ID from TenantUsers relationship
	var tenantID uuid.UUID
	if len(user.TenantUsers) > 0 {
		tenantID = user.TenantUsers[0].TenantID
	} else {
		logger.Error().
			Str("user_id", id.String()).
			Msg("User is not associated with any tenant during update")
		return nil, errors.New("user is not associated with any tenant")
	}

	// Check if email already exists (if changed and provided)
	if req.Email != nil && *req.Email != "" && *req.Email != user.Email {
		existingUser, _ := s.userRepo.GetByEmailAndTenant(c, *req.Email, tenantID)
		if existingUser != nil && existingUser.ID != id {
			logger.Warn().
				Str("email", *req.Email).
				Str("user_id", id.String()).
				Str("tenant_id", tenantID.String()).
				Msg("User update attempt with existing email")
			return nil, errors.New("email already exists")
		}
	}

	// Handle role update if provided
	if req.RoleID != nil {
		_, err := s.roleRepo.GetByID(c, *req.RoleID)
		if err != nil {
			logger.Error().
				Err(err).
				Str("role_id", req.RoleID.String()).
				Str("user_id", id.String()).
				Msg("Invalid role ID provided during user update")
			return nil, errors.New("invalid role ID")
		}

		// Get tenant user
		tenantUser, err := s.tenantUserRepo.GetByTenantAndUser(c, tenantID, user.ID)
		if err != nil {
			logger.Error().
				Err(err).
				Str("user_id", id.String()).
				Str("tenant_id", tenantID.String()).
				Msg("Tenant user not found during role update")
			return nil, errors.New("tenant user not found")
		}

		// Delete existing tenant user roles and create new one
		err = s.tenantUserRoleRepo.DeleteAllTenantUserRoles(c, tenantUser.ID)
		if err != nil {
			logger.Error().
				Err(err).
				Str("tenant_user_id", tenantUser.ID.String()).
				Msg("Failed to delete existing tenant user roles during update")
			return nil, errors.New("failed to update tenant user role")
		}

		tenantUserRole := &model.TenantUserRole{
			TenantUserID: tenantUser.ID,
			RoleID:       *req.RoleID,
		}

		err = s.tenantUserRoleRepo.Create(c, tenantUserRole)
		if err != nil {
			logger.Error().
				Err(err).
				Str("tenant_user_id", tenantUser.ID.String()).
				Str("role_id", req.RoleID.String()).
				Msg("Failed to create new tenant user role during update")
			return nil, errors.New("failed to create new tenant user role")
		}
	}

	// Update fields
	if req.Email != nil && *req.Email != "" {
		user.Email = *req.Email
	}
	if req.FullName != nil && *req.FullName != "" {
		user.FullName = *req.FullName
	}
	if req.Birthplace != nil {
		user.Birthplace = req.Birthplace
	}
	if req.Birthday != nil {
		user.Birthday = req.Birthday
	}
	if req.Gender != nil {
		user.Gender = (*model.Gender)(req.Gender)
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Address != nil {
		user.Address = req.Address
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive

		// Also update the tenant-user active status
		if len(user.TenantUsers) > 0 {
			tenantUser := &user.TenantUsers[0]
			tenantUser.IsActive = *req.IsActive
			err = s.tenantUserRepo.Update(c, tenantUser)
			if err != nil {
				logger.Error().
					Err(err).
					Str("user_id", id.String()).
					Str("tenant_id", tenantID.String()).
					Msg("Failed to update tenant-user status")
				return nil, errors.New("failed to update tenant-user status")
			}
		}
	}

	err = s.userRepo.Update(c, user)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Failed to update user in database")
		return nil, errors.New("failed to update user")
	}

	return user, nil
}

func (s *userService) Delete(c context.Context, id uuid.UUID) error {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Check if user exists
	_, err := s.userRepo.GetByID(c, id)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("User not found during delete")
		return err
	}

	err = s.userRepo.Delete(c, id)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Failed to delete user from database")
		return err
	}

	return nil
}

func (s *userService) BulkDelete(c context.Context, tenantID uuid.UUID, ids []uuid.UUID) error {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	if len(ids) == 0 {
		return errors.New("no user IDs provided for bulk delete")
	}

	// Get users that belong to the tenant to validate they exist and log properly
	users, _, err := s.userRepo.GetUsersByTenant(c, tenantID, 0, len(ids)*2, "")
	if err != nil {
		logger.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Interface("user_ids", ids).
			Msg("Failed to validate users for bulk delete")
		return errors.New("failed to validate users for bulk delete")
	}

	// Create a set of valid user IDs that belong to the tenant
	validUserMap := make(map[uuid.UUID]bool)
	for _, user := range users {
		validUserMap[user.ID] = true
	}

	// Filter IDs to only include users that belong to the tenant
	var validIDs []uuid.UUID
	var invalidIDs []uuid.UUID
	for _, id := range ids {
		if validUserMap[id] {
			validIDs = append(validIDs, id)
		} else {
			invalidIDs = append(invalidIDs, id)
		}
	}

	if len(invalidIDs) > 0 {
		logger.Warn().
			Str("tenant_id", tenantID.String()).
			Interface("invalid_ids", invalidIDs).
			Msg("Some user IDs do not belong to the tenant or do not exist")
	}

	if len(validIDs) == 0 {
		return errors.New("no valid user IDs found for bulk delete in this tenant")
	}

	// Perform bulk delete
	err = s.userRepo.BulkDelete(c, validIDs)
	if err != nil {
		logger.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Interface("user_ids", validIDs).
			Msg("Failed to bulk delete users from database")
		return errors.New("failed to bulk delete users")
	}

	return nil
}

func (s *userService) List(c context.Context, tenantID uuid.UUID, params dto.UserQueryParams) ([]model.User, *dto.PaginationMeta, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	var users []model.User
	var total int64
	var err error

	if params.RoleID != nil {
		users, total, err = s.userRepo.GetByRole(c, tenantID, *params.RoleID, offset, params.Limit)
		if err != nil {
			logger.Error().
				Err(err).
				Str("tenant_id", tenantID.String()).
				Str("role_id", params.RoleID.String()).
				Interface("params", params).
				Msg("Failed to get users by role")
		}
	} else {
		users, total, err = s.userRepo.GetUsersByTenant(c, tenantID, offset, params.Limit, params.Search)
		if err != nil {
			logger.Error().
				Err(err).
				Str("tenant_id", tenantID.String()).
				Interface("params", params).
				Msg("Failed to get users by tenant")
		}
	}

	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	meta := &dto.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalRows:  total,
		TotalPages: totalPages,
	}

	return users, meta, nil
}
