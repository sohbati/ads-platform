package container

import (
	"log"

	cachestore "cache-service/internal/cachestore"
	"cache-service/internal/simplecache/attrschema/client"
	"cache-service/internal/simplecache/attrschema/handler"
	serviceimpl "cache-service/internal/simplecache/attrschema/service/impl"
)

type AttrSchemaCacheContainer struct {
	AttrSchemaCacheHandler *handler.AttrSchemaCacheHandler
	NameToSchemaStore      *cachestore.Cache[string, string]
}

func NewAttrSchemaCacheContainer(cdnBaseURL string) *AttrSchemaCacheContainer {
	nameToSchema := cachestore.New[string, string](
		cachestore.WithTTL(cachestore.NoExpiration),
	)

	cdnClient := client.NewCDNAttrSchemaClient(cdnBaseURL)
	svc := serviceimpl.NewAttrSchemaCacheService(nameToSchema, cdnClient)
	h := handler.NewAttrSchemaCacheHandler(svc)
	log.Printf("Attr schema cache configured (cdn=%s)", cdnBaseURL)

	return &AttrSchemaCacheContainer{
		AttrSchemaCacheHandler: h,
		NameToSchemaStore:      nameToSchema,
	}
}
