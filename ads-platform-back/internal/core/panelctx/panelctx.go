// Package panelctx reads the acting staff member on /api/v1/panel requests.
//
// The CRM authenticates staff and forwards their ID in StaffIDHeader, overwriting anything the
// browser sent. Like the rest of /panel, this trusts the CRM, so the backend must not be reachable
// from outside the private network.
package panelctx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const StaffIDHeader = "X-Staff-Id"

// StaffID returns the acting staff member's ID, or false when the header is missing or invalid.
func StaffID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.GetHeader(StaffIDHeader), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
