package impl

import (
	"context"
	"errors"
	"regexp"
	"strings"

	customermodel "ads-platform/internal/business/customer/model"
	"ads-platform/internal/business/customerprofile/errorcode"
	"ads-platform/internal/business/customerprofile/model"
	"ads-platform/internal/business/customerprofile/repository"
	"ads-platform/internal/business/customerprofile/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

const maxLocationSlugs = 40

var slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,78}$`)

type profileService struct {
	profiles  repository.CustomerProfileRepository
	customers customerLookup
}

type customerLookup interface {
	GetCustomerByID(ctx context.Context, id int64) (*customermodel.Customer, error)
}

func NewCustomerProfileService(profiles repository.CustomerProfileRepository, customers customerLookup) service.CustomerProfileService {
	return &profileService{profiles: profiles, customers: customers}
}

func (s *profileService) Get(ctx context.Context, customerID int64) (*model.CustomerProfile, error) {
	if err := s.ensureCustomer(ctx, customerID); err != nil {
		return nil, err
	}

	profile, err := s.profiles.GetByCustomerID(ctx, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return emptyProfile(customerID), nil
		}
		return nil, err
	}
	if profile.LocationSlugs == nil {
		profile.LocationSlugs = []string{}
	}
	return profile, nil
}

func (s *profileService) Put(ctx context.Context, customerID int64, locationSlugs []string) (*model.CustomerProfile, error) {
	if err := s.ensureCustomer(ctx, customerID); err != nil {
		return nil, err
	}

	slugs, err := normalizeSlugs(locationSlugs)
	if err != nil {
		return nil, err
	}

	profile := &model.CustomerProfile{
		CustomerID:        customerID,
		LocationSlugs: slugs,
	}
	if err := s.profiles.Upsert(ctx, profile); err != nil {
		return nil, err
	}
	return s.Get(ctx, customerID)
}

func (s *profileService) ensureCustomer(ctx context.Context, customerID int64) error {
	if customerID <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidCustomer.Code, errorcode.ErrInvalidCustomer.HttpStatus)
	}
	_, err := s.customers.GetCustomerByID(ctx, customerID)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.NewAppError(errorcode.ErrCustomerNotFound.Code, errorcode.ErrCustomerNotFound.HttpStatus)
	}
	return err
}

func emptyProfile(customerID int64) *model.CustomerProfile {
	return &model.CustomerProfile{CustomerID: customerID, LocationSlugs: []string{}}
}

func normalizeSlugs(in []string) ([]string, error) {
	if len(in) > maxLocationSlugs {
		return nil, exception.NewAppError(errorcode.ErrTooManyLocations.Code, errorcode.ErrTooManyLocations.HttpStatus)
	}
	out := make([]string, 0, len(in))
	seen := make(map[string]struct{}, len(in))
	for _, raw := range in {
		slug := strings.ToLower(strings.TrimSpace(raw))
		if slug == "" {
			continue
		}
		if !slugPattern.MatchString(slug) {
			return nil, exception.NewAppError(errorcode.ErrInvalidLocation.Code, errorcode.ErrInvalidLocation.HttpStatus, slug)
		}
		if _, ok := seen[slug]; ok {
			continue
		}
		seen[slug] = struct{}{}
		out = append(out, slug)
	}
	if out == nil {
		out = []string{}
	}
	return out, nil
}
