package service

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/model"
	"github.com/protocyber/kelasgo-api/internal/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
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
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userService) Create(tenantID uuid.UUID, req dto.CreateUserRequest) (*model.User, error) {
	// Check if username already exists within tenant
	existingUser, _ := s.userRepo.GetByUsernameAndTenant(req.Username, tenantID)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists within tenant (if provided)
	if req.Email != "" {
		existingUser, _ = s.userRepo.GetByEmailAndTenant(req.Email, tenantID)
		if existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	// Validate role if provided
	if req.RoleID != nil {
		_, err := s.roleRepo.GetByID(*req.RoleID)
		if err != nil {
			return nil, errors.New("invalid role ID")
		}
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		TenantID:     tenantID,
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

func (s *userService) GetByID(id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) Update(id uuid.UUID, req dto.UpdateUserRequest) (*model.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email already exists (if changed and provided)
	if req.Email != nil && *req.Email != "" && *req.Email != user.Email {
		existingUser, _ := s.userRepo.GetByEmailAndTenant(*req.Email, user.TenantID)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
	}

	// Validate role if provided
	if req.RoleID != nil {
		_, err := s.roleRepo.GetByID(*req.RoleID)
		if err != nil {
			return nil, errors.New("invalid role ID")
		}
		user.RoleID = req.RoleID
	}

	// Update fields
	if req.Email != nil && *req.Email != "" {
		user.Email = *req.Email
	}
	if req.FullName != nil && *req.FullName != "" {
		user.FullName = *req.FullName
	}
	if req.Gender != nil {
		user.Gender = req.Gender
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
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, errors.New("failed to update user")
	}

	return user, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
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
	} else {
		users, total, err = s.userRepo.List(tenantID, offset, params.Limit, params.Search)
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
