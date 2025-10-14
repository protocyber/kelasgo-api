package service

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/model"
	"github.com/protocyber/kelasgo-api/internal/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
)

// UserService interface defines user service methods
type UserService interface {
	Create(tenantID uuid.UUID, req dto.CreateUserRequest) (*model.User, error)
	GetByID(id uuid.UUID) (*model.User, error)
	Update(id uuid.UUID, req dto.UpdateUserRequest) (*model.User, error)
	Delete(id uuid.UUID) error
	List(tenantID uuid.UUID, params dto.UserQueryParams) ([]model.User, *dto.PaginationMeta, error)
}

// userService implements UserService
type userService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	tenantUserRepo repository.TenantUserRepository
	userRoleRepo   repository.UserRoleRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tenantUserRepo repository.TenantUserRepository,
	userRoleRepo repository.UserRoleRepository,
) UserService {
	return &userService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		tenantUserRepo: tenantUserRepo,
		userRoleRepo:   userRoleRepo,
	}
}

func (s *userService) Create(tenantID uuid.UUID, req dto.CreateUserRequest) (*model.User, error) {
	// Check if username already exists within tenant
	existingUser, _ := s.userRepo.GetByUsernameAndTenant(req.Username, tenantID)
	if existingUser != nil {
		log.Warn().
			Str("username", req.Username).
			Str("tenant_id", tenantID.String()).
			Msg("User creation attempt with existing username")
		return nil, errors.New("username already exists")
	}

	// Check if email already exists within tenant (if provided)
	if req.Email != "" {
		existingUser, _ = s.userRepo.GetByEmailAndTenant(req.Email, tenantID)
		if existingUser != nil {
			log.Warn().
				Str("email", req.Email).
				Str("tenant_id", tenantID.String()).
				Msg("User creation attempt with existing email")
			return nil, errors.New("email already exists")
		}
	}

	// Validate role if provided
	if req.RoleID != nil {
		_, err := s.roleRepo.GetByID(*req.RoleID)
		if err != nil {
			log.Error().
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
		log.Error().
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

	err = s.userRepo.Create(user)
	if err != nil {
		log.Error().
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

	err = s.tenantUserRepo.Create(tenantUser)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.ID.String()).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to create tenant-user relationship")
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
			log.Error().
				Err(err).
				Str("user_id", user.ID.String()).
				Str("role_id", req.RoleID.String()).
				Msg("Failed to create user-role relationship")
			// If user-role creation fails, cleanup user and tenant-user
			s.tenantUserRepo.Delete(tenantUser.ID)
			s.userRepo.Delete(user.ID)
			return nil, errors.New("failed to create user-role relationship")
		}
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("username", user.Username).
		Str("tenant_id", tenantID.String()).
		Msg("User created successfully")

	return user, nil
}

func (s *userService) GetByID(id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Failed to get user by ID")
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(id uuid.UUID, req dto.UpdateUserRequest) (*model.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		log.Error().
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
		log.Error().
			Str("user_id", id.String()).
			Msg("User is not associated with any tenant during update")
		return nil, errors.New("user is not associated with any tenant")
	}

	// Check if email already exists (if changed and provided)
	if req.Email != nil && *req.Email != "" && *req.Email != user.Email {
		existingUser, _ := s.userRepo.GetByEmailAndTenant(*req.Email, tenantID)
		if existingUser != nil && existingUser.ID != id {
			log.Warn().
				Str("email", *req.Email).
				Str("user_id", id.String()).
				Str("tenant_id", tenantID.String()).
				Msg("User update attempt with existing email")
			return nil, errors.New("email already exists")
		}
	}

	// Handle role update if provided
	if req.RoleID != nil {
		_, err := s.roleRepo.GetByID(*req.RoleID)
		if err != nil {
			log.Error().
				Err(err).
				Str("role_id", req.RoleID.String()).
				Str("user_id", id.String()).
				Msg("Invalid role ID provided during user update")
			return nil, errors.New("invalid role ID")
		}

		// Delete existing user roles and create new one
		err = s.userRoleRepo.DeleteAllUserRoles(user.ID)
		if err != nil {
			log.Error().
				Err(err).
				Str("user_id", id.String()).
				Msg("Failed to delete existing user roles during update")
			return nil, errors.New("failed to update user role")
		}

		userRole := &model.UserRole{
			UserID: user.ID,
			RoleID: *req.RoleID,
		}

		err = s.userRoleRepo.Create(userRole)
		if err != nil {
			log.Error().
				Err(err).
				Str("user_id", id.String()).
				Str("role_id", req.RoleID.String()).
				Msg("Failed to create new user role during update")
			return nil, errors.New("failed to create new user role")
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
			err = s.tenantUserRepo.Update(tenantUser)
			if err != nil {
				log.Error().
					Err(err).
					Str("user_id", id.String()).
					Str("tenant_id", tenantID.String()).
					Msg("Failed to update tenant-user status")
				return nil, errors.New("failed to update tenant-user status")
			}
		}
	}

	err = s.userRepo.Update(user)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Failed to update user in database")
		return nil, errors.New("failed to update user")
	}

	log.Info().
		Str("user_id", id.String()).
		Str("username", user.Username).
		Msg("User updated successfully")

	return user, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("User not found during delete")
		return err
	}

	err = s.userRepo.Delete(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Failed to delete user from database")
		return err
	}

	log.Info().
		Str("user_id", id.String()).
		Str("username", user.Username).
		Msg("User deleted successfully")

	return nil
}

func (s *userService) List(tenantID uuid.UUID, params dto.UserQueryParams) ([]model.User, *dto.PaginationMeta, error) {
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
		users, total, err = s.userRepo.GetByRole(tenantID, *params.RoleID, offset, params.Limit)
		if err != nil {
			log.Error().
				Err(err).
				Str("tenant_id", tenantID.String()).
				Str("role_id", params.RoleID.String()).
				Int("page", params.Page).
				Int("limit", params.Limit).
				Msg("Failed to get users by role")
		}
	} else {
		users, total, err = s.userRepo.GetUsersByTenant(tenantID, offset, params.Limit, params.Search)
		if err != nil {
			log.Error().
				Err(err).
				Str("tenant_id", tenantID.String()).
				Str("search", params.Search).
				Int("page", params.Page).
				Int("limit", params.Limit).
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
