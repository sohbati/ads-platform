package container

import (
	attrschemaContainer "ads-platform-cdn/internal/business/attrschema/container"
	categoryContainer "ads-platform-cdn/internal/business/category/container"
	cityContainer "ads-platform-cdn/internal/business/city/container"
	"ads-platform-cdn/internal/core/config"
)

type AppContainer struct {
	Config     *config.Config
	Category   *categoryContainer.CategoryContainer
	City       *cityContainer.CityContainer
	AttrSchema *attrschemaContainer.AttrSchemaContainer
}

func NewAppContainer(cfg *config.Config) (*AppContainer, error) {
	return &AppContainer{
		Config:     cfg,
		Category:   categoryContainer.NewCategoryContainer(cfg.CategoryJSON),
		City:       cityContainer.NewCityContainer(cfg.CitiesJSON),
		AttrSchema: attrschemaContainer.NewAttrSchemaContainer(cfg.AttrSchemasJSON),
	}, nil
}
