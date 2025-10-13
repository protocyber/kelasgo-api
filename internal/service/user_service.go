package service

import (
	"errors"
	"math"

	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/model"
	"github.com/protocyber/kelasgo-api/internal/repository"
)

// UserService interface defines user service methods
type UserService interface {
	Create(req dto.CreateUserRequest) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	Update(id uint, req dto.UpdateUserRequest) (*model.User, error)
	Delete(id uint) error
	List(params dto.UserQueryParams) ([]model.User, *dto.PaginationMeta, error)
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

func (s *userService) Create(req dto.CreateUserRequest) (*model.User, error) {
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

	// Validate role if provided
	if req.RoleID.Valid {
		_, err := s.roleRepo.GetByID(uint(req.RoleID.Int64))
		if err != nil {
			return nil, errors.New("invalid role ID")
		}
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
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

func (s *userService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) Update(id uint, req dto.UpdateUserRequest) (*model.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email already exists (if changed and provided)
	if req.Email != "" && req.Email != user.Email {
		existingUser, _ := s.userRepo.GetByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
	}

	// Validate role if provided
	if req.RoleID.Valid {
		_, err := s.roleRepo.GetByID(uint(req.RoleID.Int64))
		if err != nil {
			return nil, errors.New("invalid role ID")
		}
		user.RoleID = req.RoleID
	}

	// Update fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.DateOfBirth.Valid {
		user.DateOfBirth = req.DateOfBirth
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Address != "" {
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

func (s *userService) Delete(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}

func (s *userService) List(params dto.UserQueryParams) ([]model.User, *dto.PaginationMeta, error) {
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

	if params.RoleID > 0 {
		users, total, err = s.userRepo.GetByRole(uint(params.RoleID), offset, params.Limit)
	} else {
		users, total, err = s.userRepo.List(offset, params.Limit, params.Search)
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

// Helper function to hash password (placeholder implementation)
func hashPassword(password string) (string, error) {
	// This should use util.HashPassword but to avoid circular import
	// we'll implement it here or refactor later
	return password + "_hashed", nil // Placeholder
}
