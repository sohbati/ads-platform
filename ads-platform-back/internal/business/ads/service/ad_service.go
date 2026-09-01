package service

import (
	"context"
	"encoding/json"
	"io"

	"ads-platform/internal/business/ads/model"
)

type CreateAdInput struct {
	UserID       int64
	CategoryID   int
	CityID       int
	Title        string
	Description  string
	Latitude     *float64
	Longitude    *float64
	Neighborhood string
	PriceAmount  *int64
	PriceType    string
	Currency     string
	Attrs        json.RawMessage
	Contact      json.RawMessage
	Pictures     []PictureInput
	KeepMedia    *[]string
}

type PictureInput struct {
	Filename    string
	ContentType string
	Size        int64
	Body        io.Reader
}

type AdService interface {
	Create(ctx context.Context, in CreateAdInput) (*model.Ad, error)
	GetPublic(ctx context.Context, adID int64) (*model.PublicAd, error)
	GetPublicContact(ctx context.Context, adID int64) (*model.PublicContact, error)
	GetForOwner(ctx context.Context, userID, adID int64) (*model.Ad, error)
	Update(ctx context.Context, adID int64, in CreateAdInput) (*model.Ad, error)
	ListByUser(ctx context.Context, userID int64) ([]model.UserAdItem, error)
	ListStats(ctx context.Context, userID int64, from, to string) (*model.AdStatsResponse, error)
}
