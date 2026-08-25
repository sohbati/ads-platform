package container

import (
	"ads-platform-cdn/internal/business/attrschema/handler"
	serviceimpl "ads-platform-cdn/internal/business/attrschema/service/impl"
)

type AttrSchemaContainer struct {
	APIHandler *handler.APIHandler
}

func NewAttrSchemaContainer(attrSchemasJSONPath string) *AttrSchemaContainer {
	svc := serviceimpl.NewAttrSchemaService(attrSchemasJSONPath)
	return &AttrSchemaContainer{
		APIHandler: handler.NewAPIHandler(svc),
	}
}
