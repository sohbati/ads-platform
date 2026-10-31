package impl

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"ads-platform/internal/business/setting/errorcode"
	"ads-platform/internal/business/setting/model"
	"ads-platform/internal/business/setting/repository"
	"ads-platform/internal/business/setting/service"
	"ads-platform/internal/core/exception"
)

type definition struct {
	kind     string
	min, max int
	fallback string
}

// registry lists the settings the panel may read and change, with their validation.
var registry = map[string]definition{
	model.AdsPendingMaxDays: {kind: "int", min: 1, max: 365, fallback: "7"},
}

type settingService struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) service.SettingService {
	return &settingService{repo: repo}
}

func (s *settingService) List(ctx context.Context) ([]model.SettingView, error) {
	saved, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	byKey := make(map[string]model.AppSetting, len(saved))
	for _, item := range saved {
		byKey[item.Key] = item
	}

	out := make([]model.SettingView, 0, len(registry))
	for _, key := range sortedKeys() {
		item, ok := byKey[key]
		out = append(out, view(key, registry[key], &item, ok))
	}
	return out, nil
}

func (s *settingService) Update(ctx context.Context, key, raw string, staffID int64) (*model.SettingView, error) {
	if staffID <= 0 {
		return nil, exception.NewAppError(errorcode.ErrStaffRequired.Code, errorcode.ErrStaffRequired.HttpStatus)
	}
	def, ok := registry[key]
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrUnknownKey.Code, errorcode.ErrUnknownKey.HttpStatus, key)
	}
	value, ok := normalize(def, raw)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidValue.Code, errorcode.ErrInvalidValue.HttpStatus, key)
	}

	setting := &model.AppSetting{Key: key, Value: value, UpdatedByStaffID: &staffID, UpdatedAt: time.Now().UTC()}
	if err := s.repo.Upsert(ctx, setting); err != nil {
		return nil, err
	}
	result := view(key, def, setting, true)
	return &result, nil
}

func (s *settingService) Int(ctx context.Context, key string) (int, error) {
	def, ok := registry[key]
	if !ok || def.kind != "int" {
		return 0, exception.NewAppError(errorcode.ErrUnknownKey.Code, errorcode.ErrUnknownKey.HttpStatus, key)
	}
	saved, err := s.repo.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	raw := def.fallback
	if saved != nil {
		if value, ok := normalize(def, saved.Value); ok {
			raw = value
		}
	}
	return strconv.Atoi(raw)
}

func normalize(def definition, raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	switch def.kind {
	case "int":
		n, err := strconv.Atoi(value)
		if err != nil || n < def.min || n > def.max {
			return "", false
		}
		return strconv.Itoa(n), true
	default:
		return "", false
	}
}

func view(key string, def definition, saved *model.AppSetting, found bool) model.SettingView {
	min, max := def.min, def.max
	out := model.SettingView{Key: key, Value: def.fallback, Type: def.kind, Min: &min, Max: &max, Default: def.fallback}
	if found {
		out.Value = saved.Value
		out.UpdatedByStaffID = saved.UpdatedByStaffID
		updated := saved.UpdatedAt
		out.UpdatedAt = &updated
	}
	return out
}

func sortedKeys() []string {
	keys := make([]string, 0, len(registry))
	for key := range registry {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
