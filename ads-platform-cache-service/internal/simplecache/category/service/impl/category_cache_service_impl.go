package impl

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"sync"

	"cache-service/internal/cachestore"
	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/category/client"
	"cache-service/internal/simplecache/category/errorcode"
	"cache-service/internal/simplecache/category/model"
	"cache-service/internal/simplecache/category/service"
)

type categoryCacheService struct {
	slugToCategory *cachestore.Cache[string, string] // slug → category JSON
	idToSlug       *cachestore.Cache[int, string]    // id → slug
	cdn            client.CDNCategoryClient
	loadMu         sync.Mutex
}

func NewCategoryCacheService(
	slugToCategory *cachestore.Cache[string, string],
	idToSlug *cachestore.Cache[int, string],
	cdn client.CDNCategoryClient,
) service.CategoryCacheService {
	return &categoryCacheService{
		slugToCategory: slugToCategory,
		idToSlug:       idToSlug,
		cdn:            cdn,
	}
}

func (s *categoryCacheService) GetBySlugs(ctx context.Context, slugs []string, includeDescendants bool) ([]model.Category, error) {
	clean := make([]string, 0, len(slugs))
	for _, slug := range slugs {
		slug = strings.TrimSpace(slug)
		if slug != "" {
			clean = append(clean, slug)
		}
	}
	if len(clean) == 0 {
		return nil, exception.NewAppError(errorcode.ErrSlugsEmpty.Code, errorcode.ErrSlugsEmpty.HttpStatus)
	}
	if err := s.ensureLoaded(ctx); err != nil {
		return nil, err
	}

	out := make([]model.Category, 0, len(clean))
	for _, slug := range clean {
		category, ok := s.categoryBySlug(ctx, slug)
		if !ok {
			continue
		}
		if includeDescendants {
			category.DescendantIDs = s.descendantIDs(ctx, category.ID)
		}
		out = append(out, category)
	}
	return out, nil
}

func (s *categoryCacheService) GetByIDs(ctx context.Context, ids []int, includeDescendants bool) ([]model.Category, error) {
	if len(ids) == 0 {
		return nil, exception.NewAppError(errorcode.ErrIDsEmpty.Code, errorcode.ErrIDsEmpty.HttpStatus)
	}
	if err := s.ensureLoaded(ctx); err != nil {
		return nil, err
	}

	out := make([]model.Category, 0, len(ids))
	for _, id := range ids {
		slug, err := s.idToSlug.Get(ctx, id)
		if err != nil {
			continue
		}
		category, ok := s.categoryBySlug(ctx, slug)
		if !ok {
			continue
		}
		if includeDescendants {
			category.DescendantIDs = s.descendantIDs(ctx, category.ID)
		}
		out = append(out, category)
	}
	return out, nil
}

// descendantIDs returns self + all descendant category ids, derived from the
// Path field (comma-separated ancestor ids root→self).
func (s *categoryCacheService) descendantIDs(ctx context.Context, id int) []int {
	ids := make([]int, 0, 8)
	for _, slug := range s.slugToCategory.Keys() {
		category, ok := s.categoryBySlug(ctx, slug)
		if !ok {
			continue
		}
		if pathContains(category.Path, id) || category.ID == id {
			ids = append(ids, category.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func pathContains(path string, id int) bool {
	if path == "" {
		return false
	}
	for _, part := range strings.Split(path, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil && v == id {
			return true
		}
	}
	return false
}

func (s *categoryCacheService) categoryBySlug(ctx context.Context, slug string) (model.Category, bool) {
	raw, err := s.slugToCategory.Get(ctx, slug)
	if err != nil {
		return model.Category{}, false
	}
	var category model.Category
	if err := json.Unmarshal([]byte(raw), &category); err != nil {
		return model.Category{}, false
	}
	return category, true
}

func (s *categoryCacheService) ensureLoaded(ctx context.Context) error {
	if s.slugToCategory.Count() > 0 && s.idToSlug.Count() > 0 {
		return nil
	}

	s.loadMu.Lock()
	defer s.loadMu.Unlock()

	if s.slugToCategory.Count() > 0 && s.idToSlug.Count() > 0 {
		return nil
	}

	categories, err := s.cdn.ListCategories(ctx)
	if err != nil {
		return exception.NewAppError(
			errorcode.ErrCDNUnavailable.Code,
			errorcode.ErrCDNUnavailable.HttpStatus,
		).WithCause(err)
	}
	if len(categories) == 0 {
		return exception.NewAppError(errorcode.ErrCacheEmpty.Code, errorcode.ErrCacheEmpty.HttpStatus)
	}

	for _, category := range categories {
		if category.Slug == "" {
			continue
		}
		raw, err := json.Marshal(category)
		if err != nil {
			continue
		}
		_ = s.slugToCategory.Set(ctx, category.Slug, string(raw))
		_ = s.idToSlug.Set(ctx, category.ID, category.Slug)
	}
	return nil
}
