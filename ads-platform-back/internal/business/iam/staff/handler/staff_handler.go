package handler

import (
	"ads-platform/internal/business/iam/staff/service"
)

// StaffHandler is reserved for future panel staff listing. IAM lookups use StaffRepository.
type StaffHandler struct {
	staffService service.StaffService
}

func NewStaffHandler(staffService service.StaffService) *StaffHandler {
	return &StaffHandler{staffService: staffService}
}
