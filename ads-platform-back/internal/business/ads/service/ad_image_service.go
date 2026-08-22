package service

import (
	"context"

	"ads-platform/internal/business/ads/model"
)

type AdImageService interface {
	// Register validates the metadata, generates a unique object key, and
	// stores a "pending" record. The caller then uploads the file to object
	// storage under that key.
	Register(ctx context.Context, userID int64, originalFilename, contentType string, fileSize int64) (*model.AdImage, error)

	// MarkUploaded confirms the file landed in object storage and records
	// its checksum and final size.
	MarkUploaded(ctx context.Context, userID, imageID int64, checksum string, fileSize int64) (*model.AdImage, error)

	// Get returns the image metadata, scoped to its owner.
	Get(ctx context.Context, userID, imageID int64) (*model.AdImage, error)

	// Delete soft-deletes the image so a cleanup job can purge the file later.
	Delete(ctx context.Context, userID, imageID int64) error
}
