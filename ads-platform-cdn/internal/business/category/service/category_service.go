package service

import "ads-platform-cdn/internal/business/category/model"

type CategoryService interface {
	List() ([]model.Category, error)
}
