package impl

import (
	"context"

	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/ads/repository"

	"gorm.io/gorm"
)

type adsAttrsJSONSchemaTemplateRepository struct {
	db *gorm.DB
}

func NewAdsAttrsJSONSchemaTemplateRepository(db *gorm.DB) repository.AdsAttrsJSONSchemaTemplateRepository {
	return &adsAttrsJSONSchemaTemplateRepository{db: db}
}

func (r *adsAttrsJSONSchemaTemplateRepository) Create(ctx context.Context, template *model.AdsAttrsJSONSchemaTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *adsAttrsJSONSchemaTemplateRepository) GetByID(ctx context.Context, id int64) (*model.AdsAttrsJSONSchemaTemplate, error) {
	var template model.AdsAttrsJSONSchemaTemplate
	if err := r.db.WithContext(ctx).First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *adsAttrsJSONSchemaTemplateRepository) GetByName(ctx context.Context, name string) (*model.AdsAttrsJSONSchemaTemplate, error) {
	var template model.AdsAttrsJSONSchemaTemplate
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *adsAttrsJSONSchemaTemplateRepository) Update(ctx context.Context, template *model.AdsAttrsJSONSchemaTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}
