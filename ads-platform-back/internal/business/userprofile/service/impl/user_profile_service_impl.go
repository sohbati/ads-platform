package impl

import (
	"context"
	"errors"
	"regexp"
	"strings"

	usermodel "ads-platform/internal/business/user/model"
	userrepo "ads-platform/internal/business/user/repository"
	"ads-platform/internal/business/userprofile/errorcode"
	"ads-platform/internal/business/userprofile/model"
	"ads-platform/internal/business/userprofile/repository"
	"ads-platform/internal/business/userprofile/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

const maxLocationSlugs = 40

var slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,78}$`)

type profileService struct {
	profiles repository.UserProfileRepository
	users    userLookup
}

type userLookup interface {
	GetUserByID(ctx context.Context, id int64) (*usermodel.User, error)
}

func NewUserProfileService(profiles repository.UserProfileRepository, users userrepo.UserRepository) service.UserProfileService {
	return &profileService{profiles: profiles, users: users}
}

func (s *profileService) Get(ctx context.Context, userID int64) (*model.UserProfile, error) {
	if err := s.ensureUser(ctx, userID); err != nil {
		return nil, err
	}

	profile, err := s.profiles.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return emptyProfile(userID), nil
		}
		return nil, err
	}
	if profile.LocationSlugs == nil {
		profile.LocationSlugs = []string{}
	}
	return profile, nil
}

func (s *profileService) Put(ctx context.Context, userID int64, locationSlugs []string) (*model.UserProfile, error) {
	if err := s.ensureUser(ctx, userID); err != nil {
		return nil, err
	}

	slugs, err := normalizeSlugs(locationSlugs)
	if err != nil {
		return nil, err
	}

	profile := &model.UserProfile{
		UserID:        userID,
		LocationSlugs: slugs,
	}
	if err := s.profiles.Upsert(ctx, profile); err != nil {
		return nil, err
	}
	return s.Get(ctx, userID)
}

func (s *profileService) ensureUser(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidUser.Code, errorcode.ErrInvalidUser.HttpStatus)
	}
	_, err := s.users.GetUserByID(ctx, userID)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.NewAppError(errorcode.ErrUserNotFound.Code, errorcode.ErrUserNotFound.HttpStatus)
	}
	return err
}

func emptyProfile(userID int64) *model.UserProfile {
	return &model.UserProfile{UserID: userID, LocationSlugs: []string{}}
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
