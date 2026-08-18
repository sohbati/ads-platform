package impl

import (
	"context"
	"encoding/json"
	"sync"

	"cache-service/internal/cachestore"
	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/city/client"
	"cache-service/internal/simplecache/city/errorcode"
	"cache-service/internal/simplecache/city/model"
	"cache-service/internal/simplecache/city/service"
)

type cityCacheService struct {
	idToCity *cachestore.Cache[int, string]  // id → city JSON
	slugToID *cachestore.Cache[string, int]  // slug → id
	cdn      client.CDNCityClient
	loadMu   sync.Mutex
}

func NewCityCacheService(
	idToCity *cachestore.Cache[int, string],
	slugToID *cachestore.Cache[string, int],
	cdn client.CDNCityClient,
) service.CityCacheService {
	return &cityCacheService{
		idToCity: idToCity,
		slugToID: slugToID,
		cdn:      cdn,
	}
}

func (s *cityCacheService) GetByIDs(ctx context.Context, ids []int) ([]model.City, error) {
	if len(ids) == 0 {
		return nil, exception.NewAppError(errorcode.ErrIDsEmpty.Code, errorcode.ErrIDsEmpty.HttpStatus)
	}
	if err := s.ensureLoaded(ctx); err != nil {
		return nil, err
	}

	out := make([]model.City, 0, len(ids))
	for _, id := range ids {
		raw, err := s.idToCity.Get(ctx, id)
		if err != nil {
			continue
		}
		var city model.City
		if err := json.Unmarshal([]byte(raw), &city); err != nil {
			continue
		}
		out = append(out, city)
	}
	return out, nil
}

func (s *cityCacheService) GetBySlugs(ctx context.Context, slugs []string) ([]model.City, error) {
	if len(slugs) == 0 {
		return nil, exception.NewAppError(errorcode.ErrSlugsEmpty.Code, errorcode.ErrSlugsEmpty.HttpStatus)
	}
	if err := s.ensureLoaded(ctx); err != nil {
		return nil, err
	}

	out := make([]model.City, 0, len(slugs))
	for _, slug := range slugs {
		id, err := s.slugToID.Get(ctx, slug)
		if err != nil {
			continue
		}
		raw, err := s.idToCity.Get(ctx, id)
		if err != nil {
			continue
		}
		var city model.City
		if err := json.Unmarshal([]byte(raw), &city); err != nil {
			continue
		}
		out = append(out, city)
	}
	return out, nil
}

func (s *cityCacheService) ensureLoaded(ctx context.Context) error {
	if s.idToCity.Count() > 0 && s.slugToID.Count() > 0 {
		return nil
	}

	s.loadMu.Lock()
	defer s.loadMu.Unlock()

	if s.idToCity.Count() > 0 && s.slugToID.Count() > 0 {
		return nil
	}

	cities, err := s.cdn.ListCities(ctx)
	if err != nil {
		return exception.NewAppError(
			errorcode.ErrCDNUnavailable.Code,
			errorcode.ErrCDNUnavailable.HttpStatus,
		).WithCause(err)
	}
	if len(cities) == 0 {
		return exception.NewAppError(errorcode.ErrCityCacheEmpty.Code, errorcode.ErrCityCacheEmpty.HttpStatus)
	}

	for _, city := range cities {
		raw, err := json.Marshal(city)
		if err != nil {
			continue
		}
		_ = s.idToCity.Set(ctx, city.ID, string(raw))
		if city.Slug != "" {
			_ = s.slugToID.Set(ctx, city.Slug, city.ID)
		}
	}
	return nil
}
