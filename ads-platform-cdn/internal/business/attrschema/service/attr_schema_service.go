package service

import "ads-platform-cdn/internal/business/attrschema/model"

type AttrSchemaService interface {
	List() ([]model.AttrSchema, error)
}
