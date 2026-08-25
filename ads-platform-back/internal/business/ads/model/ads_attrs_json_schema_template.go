package model

import (
	"encoding/json"
	"time"
)

// AdsAttrsJSONSchemaTemplate maps to ads_platform_schema.ads_attrs_json_schema_template.
// Category.adsAttrsJsonSchemaTemplateName points at Template.Name.
type AdsAttrsJSONSchemaTemplate struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name" gorm:"type:varchar(100);not null;unique"`
	Title       string          `json:"title" gorm:"type:varchar(200);not null"`
	Description *string         `json:"description"`
	JSONSchema  json.RawMessage `json:"json_schema" gorm:"column:json_schema;type:jsonb;not null"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AdsAttrsJSONSchemaTemplate) TableName() string {
	return "ads_attrs_json_schema_template"
}
