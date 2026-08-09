package impl

import (
	"ads-platform/internal/business/user/errorcode"
	"ads-platform/internal/business/user/model"
	"ads-platform/internal/business/user/repository"
	"ads-platform/internal/business/user/service"

	"ads-platform/internal/core/exception"
	"context"
	"errors"

	"gorm.io/gorm"
)

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with validation
func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	// Validate required fields

	if user.Mobile == "" {
		return errors.New("mobile_is_required")
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *userService) GetUserByMobile(ctx context.Context, mobile string) (*model.User, error) {
	if mobile == "" {
		return nil, exception.NewAppError(
			errorcode.ErrMobileEmpty.Code, errorcode.ErrMobileEmpty.HttpStatus, mobile)
	}
	user, err := s.userRepo.GetUserByMobile(ctx, mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(
				errorcode.ErrUserNotFound.Code,errorcode.ErrUserNotFound.HttpStatus, mobile).WithCause(err)
		} else {
			return nil, err
		}
	}
	return user, nil
}

func (s *userService) RegisterByMobile(ctx context.Context, mobile string) (*model.User, error) {
	if mobile == "" {
		return nil, exception.NewAppError(
			errorcode.ErrMobileEmpty.Code, errorcode.ErrMobileEmpty.HttpStatus, mobile)
	}

	existing, err := s.userRepo.GetUserByMobile(ctx, mobile)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &model.User{
		Mobile:     mobile,
		Name:       mobile,
		NationalId: "",
	}
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		existing, getErr := s.userRepo.GetUserByMobile(ctx, mobile)
		if getErr == nil {
			return existing, nil
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid_user_id")
	}
	return s.userRepo.GetUserByID(ctx, id)
}

// GetUsers retrieves all users with pagination
func (s *userService) GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.userRepo.GetUsers(ctx, limit, offset)
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {

	if user.Mobile == "" {
		return errors.New("last_name_is_required")
	}

	return s.userRepo.UpdateUser(ctx, user)
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid_user_id")
	}
	return s.userRepo.DeleteUser(ctx, id)
}

// ValidateUserData validates user data
func (s *userService) ValidateUserData(ctx context.Context, user *model.User) error {
	if user.Mobile == "" {
		return errors.New("mobile_is_required")
	}
	return nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
