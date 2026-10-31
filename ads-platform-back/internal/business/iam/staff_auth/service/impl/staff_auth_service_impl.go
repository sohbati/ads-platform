package impl

import (
	"context"
	"errors"

	actionmodel "ads-platform/internal/business/iam/action/model"
	"ads-platform/internal/business/iam/internal/password"
	staffmodel "ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/business/iam/staff_auth/errorcode"
	"ads-platform/internal/business/iam/staff_auth/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/mobile"

	"gorm.io/gorm"
)

type staffLookup interface {
	GetStaffByID(ctx context.Context, id int64) (*staffmodel.Staff, error)
	GetStaffByMobile(ctx context.Context, mobile string) (*staffmodel.Staff, error)
}

type actionLookup interface {
	ListActions(ctx context.Context, staffID int64) ([]actionmodel.Action, error)
}

type staffAuthService struct {
	staff      staffLookup
	actions    actionLookup
	mobileNorm *mobile.Normalizer
}

func NewStaffAuthService(staff staffLookup, actions actionLookup, mobileNorm *mobile.Normalizer) service.StaffAuthService {
	return &staffAuthService{staff: staff, actions: actions, mobileNorm: mobileNorm}
}

func (s *staffAuthService) Login(ctx context.Context, rawMobile, rawPassword string) (*service.StaffSession, error) {
	invalid := exception.NewAppError(errorcode.ErrInvalidCredentials.Code, errorcode.ErrInvalidCredentials.HttpStatus)

	var staff *staffmodel.Staff
	if normalized, err := s.mobileNorm.Normalize(rawMobile); err == nil {
		staff, err = s.staff.GetStaffByMobile(ctx, normalized)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	hash := ""
	if staff != nil {
		hash = staff.PasswordHash
	}
	// Always run Verify so unknown mobiles cost the same time as wrong passwords.
	if !password.Verify(hash, rawPassword) || staff == nil || !staff.IsActive {
		return nil, invalid
	}
	return s.build(ctx, staff)
}

func (s *staffAuthService) Session(ctx context.Context, staffID int64) (*service.StaffSession, error) {
	invalid := exception.NewAppError(errorcode.ErrSessionInvalid.Code, errorcode.ErrSessionInvalid.HttpStatus)
	if staffID <= 0 {
		return nil, invalid
	}
	staff, err := s.staff.GetStaffByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, invalid
		}
		return nil, err
	}
	if !staff.IsActive {
		return nil, invalid
	}
	return s.build(ctx, staff)
}

func (s *staffAuthService) build(ctx context.Context, staff *staffmodel.Staff) (*service.StaffSession, error) {
	if staff.IsSuperAdmin {
		return &service.StaffSession{Staff: *staff, Permissions: []string{service.SuperAdminPermission}}, nil
	}
	actions, err := s.actions.ListActions(ctx, staff.ID)
	if err != nil {
		return nil, err
	}
	permissions := make([]string, 0, len(actions))
	for _, action := range actions {
		permissions = append(permissions, action.Code)
	}
	return &service.StaffSession{Staff: *staff, Permissions: permissions}, nil
}
