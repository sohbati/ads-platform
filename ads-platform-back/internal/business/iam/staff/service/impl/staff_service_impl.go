package impl

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"ads-platform/internal/business/iam/internal/password"
	"ads-platform/internal/business/iam/internal/pgerr"
	"ads-platform/internal/business/iam/staff/errorcode"
	"ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/business/iam/staff/repository"
	"ads-platform/internal/business/iam/staff/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/mobile"

	"gorm.io/gorm"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

var nationalIdPattern = regexp.MustCompile(`^[0-9]{10}$`)

type staffService struct {
	staffRepo  repository.StaffRepository
	mobileNorm *mobile.Normalizer
}

func NewStaffService(staffRepo repository.StaffRepository, mobileNorm *mobile.Normalizer) service.StaffService {
	return &staffService{staffRepo: staffRepo, mobileNorm: mobileNorm}
}

func (s *staffService) Create(ctx context.Context, input service.StaffInput) (*model.Staff, error) {
	staff := &model.Staff{IsActive: true}
	if err := s.apply(staff, input); err != nil {
		return nil, err
	}
	hash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	staff.PasswordHash = hash
	if err := s.staffRepo.CreateStaff(ctx, staff); err != nil {
		return nil, mapWriteError(err)
	}
	return staff, nil
}

func (s *staffService) GetStaffByID(ctx context.Context, id int64) (*model.Staff, error) {
	if id <= 0 {
		return nil, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	staff, err := s.staffRepo.GetStaffByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return nil, err
	}
	return staff, nil
}

func (s *staffService) List(ctx context.Context, search string, page, pageSize int) (*service.StaffPage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	items, total, err := s.staffRepo.ListStaff(ctx, searchTerms(search), pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	return &service.StaffPage{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *staffService) Update(ctx context.Context, id int64, input service.StaffInput) (*model.Staff, error) {
	staff, err := s.GetStaffByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.apply(staff, input); err != nil {
		return nil, err
	}
	if err := s.staffRepo.UpdateStaff(ctx, staff); err != nil {
		return nil, mapWriteError(err)
	}
	return s.staffRepo.GetStaffByID(ctx, id)
}

func (s *staffService) SetPassword(ctx context.Context, id int64, raw string) error {
	staff, err := s.GetStaffByID(ctx, id)
	if err != nil {
		return err
	}
	hash, err := hashPassword(raw)
	if err != nil {
		return err
	}
	staff.PasswordHash = hash
	return s.staffRepo.UpdateStaff(ctx, staff)
}

func (s *staffService) UpsertSuperAdmin(ctx context.Context, input service.StaffInput) (*model.Staff, error) {
	normalizedMobile, err := s.mobileNorm.Normalize(input.Mobile)
	if err != nil {
		return nil, exception.NewAppError(errorcode.ErrInvalidMobile.Code, errorcode.ErrInvalidMobile.HttpStatus, input.Mobile)
	}
	hash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	staff, err := s.staffRepo.GetStaffByMobile(ctx, normalizedMobile)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		staff = &model.Staff{}
		if err := s.apply(staff, input); err != nil {
			return nil, err
		}
		staff.PasswordHash, staff.IsSuperAdmin, staff.IsActive = hash, true, true
		if err := s.staffRepo.CreateStaff(ctx, staff); err != nil {
			return nil, mapWriteError(err)
		}
		return staff, nil
	case err != nil:
		return nil, err
	}

	staff.PasswordHash, staff.IsSuperAdmin, staff.IsActive = hash, true, true
	if err := s.staffRepo.UpdateStaff(ctx, staff); err != nil {
		return nil, mapWriteError(err)
	}
	return staff, nil
}

func (s *staffService) Delete(ctx context.Context, id int64) error {
	staff, err := s.GetStaffByID(ctx, id)
	if err != nil {
		return err
	}
	if staff.IsSuperAdmin {
		return exception.NewAppError(errorcode.ErrSuperAdminProtected.Code, errorcode.ErrSuperAdminProtected.HttpStatus)
	}
	if err := s.staffRepo.DeleteStaff(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return err
	}
	return nil
}

// apply validates input and copies the normalized values onto staff.
func (s *staffService) apply(staff *model.Staff, input service.StaffInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" || len(name) > 100 {
		return exception.NewAppError(errorcode.ErrInvalidName.Code, errorcode.ErrInvalidName.HttpStatus)
	}
	normalizedMobile, err := s.mobileNorm.Normalize(input.Mobile)
	if err != nil {
		return exception.NewAppError(errorcode.ErrInvalidMobile.Code, errorcode.ErrInvalidMobile.HttpStatus, input.Mobile)
	}
	nationalId := mobile.ToLatinDigits(strings.TrimSpace(input.NationalId))
	if !nationalIdPattern.MatchString(nationalId) {
		return exception.NewAppError(errorcode.ErrInvalidNationalId.Code, errorcode.ErrInvalidNationalId.HttpStatus, input.NationalId)
	}
	if input.IsActive != nil {
		// Super admins are managed only through the create-admin command.
		if staff.IsSuperAdmin && !*input.IsActive {
			return exception.NewAppError(errorcode.ErrSuperAdminProtected.Code, errorcode.ErrSuperAdminProtected.HttpStatus)
		}
		staff.IsActive = *input.IsActive
	}
	staff.Name = name
	staff.Mobile = normalizedMobile
	staff.NationalId = nationalId
	return nil
}

func hashPassword(raw string) (string, error) {
	if !password.Valid(raw) {
		return "", exception.NewAppError(errorcode.ErrInvalidPassword.Code, errorcode.ErrInvalidPassword.HttpStatus)
	}
	return password.Hash(raw)
}

func mapWriteError(err error) error {
	constraint, ok := pgerr.UniqueViolationConstraint(err)
	if !ok {
		return err
	}
	if strings.Contains(constraint, "national_id") {
		return exception.NewAppError(errorcode.ErrDuplicateNationalId.Code, errorcode.ErrDuplicateNationalId.HttpStatus)
	}
	return exception.NewAppError(errorcode.ErrDuplicateMobile.Code, errorcode.ErrDuplicateMobile.HttpStatus)
}

// searchTerms also matches local mobile input ("0912…") against stored E.164 numbers ("+98912…").
func searchTerms(raw string) []string {
	term := mobile.ToLatinDigits(strings.TrimSpace(raw))
	if term == "" {
		return nil
	}
	terms := []string{term}
	if len(term) > 1 && term[0] == '0' && strings.Trim(term, "0123456789") == "" {
		terms = append(terms, term[1:])
	}
	return terms
}
