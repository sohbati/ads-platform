package model

import "time"

// Image rows are written as "uploaded" when an ad is created with pictures.
// Soft-delete (status deleted + deleted_at) lets a later job purge MinIO objects.
const (
	ImageStatusPending  = "pending"
	ImageStatusUploaded = "uploaded"
	ImageStatusDeleted  = "deleted"
)

// AdImage maps to ads_platform_schema.ad_images
type AdImage struct {
	ID               int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           int64      `json:"user_id" gorm:"not null"`
	AdID             *int64     `json:"ad_id"`
	ObjectKey        string     `json:"object_key" gorm:"type:varchar(255);not null;unique"`
	OriginalFilename string     `json:"original_filename" gorm:"type:varchar(255);not null"`
	ContentType      string     `json:"content_type" gorm:"type:varchar(100);not null"`
	FileSize         int64      `json:"file_size" gorm:"not null;default:0"`
	Status           string     `json:"status" gorm:"type:varchar(20);not null;default:pending"`
	Checksum         *string    `json:"checksum" gorm:"type:varchar(64)"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UploadedAt       *time.Time `json:"uploaded_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

func (AdImage) TableName() string {
	return "ad_images"
}
