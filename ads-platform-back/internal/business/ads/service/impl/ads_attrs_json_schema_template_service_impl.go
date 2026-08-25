package impl

import (
	"context"
	"strings"

	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/ads/repository"
	"ads-platform/internal/business/ads/service"
)

type adsAttrsJSONSchemaTemplateService struct {
	repo repository.AdsAttrsJSONSchemaTemplateRepository
}

func NewAdsAttrsJSONSchemaTemplateService(repo repository.AdsAttrsJSONSchemaTemplateRepository) service.AdsAttrsJSONSchemaTemplateService {
	return &adsAttrsJSONSchemaTemplateService{repo: repo}
}

func (s *adsAttrsJSONSchemaTemplateService) GetByName(ctx context.Context, name string) (*model.AdsAttrsJSONSchemaTemplate, error) {
	return s.repo.GetByName(ctx, strings.TrimSpace(name))
}

func (s *adsAttrsJSONSchemaTemplateService) GetByID(ctx context.Context, id int64) (*model.AdsAttrsJSONSchemaTemplate, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *adsAttrsJSONSchemaTemplateService) Create(ctx context.Context, template *model.AdsAttrsJSONSchemaTemplate) error {
	return s.repo.Create(ctx, template)
}

func (s *adsAttrsJSONSchemaTemplateService) Update(ctx context.Context, template *model.AdsAttrsJSONSchemaTemplate) error {
	return s.repo.Update(ctx, template)
}
