package container

import (
	adReviewContainer "ads-platform/internal/business/adreview/container"
	adsContainer "ads-platform/internal/business/ads/container"
	customerContainer "ads-platform/internal/business/customer/container"
	customerprofileContainer "ads-platform/internal/business/customerprofile/container"
	actionContainer "ads-platform/internal/business/iam/action/container"
	roleContainer "ads-platform/internal/business/iam/role/container"
	roleActionContainer "ads-platform/internal/business/iam/role_action/container"
	staffContainer "ads-platform/internal/business/iam/staff/container"
	staffAuthContainer "ads-platform/internal/business/iam/staff_auth/container"
	staffGroupContainer "ads-platform/internal/business/iam/staff_group/container"
	staffRoleContainer "ads-platform/internal/business/iam/staff_role/container"
	otpContainer "ads-platform/internal/business/otp/container"
	searchContainer "ads-platform/internal/business/search/container"
	settingContainer "ads-platform/internal/business/setting/container"
	"ads-platform/internal/core/config"

	"gorm.io/gorm"
)

type AppContainer struct {
	Customer        *customerContainer.CustomerContainer
	CustomerProfile *customerprofileContainer.CustomerProfileContainer
	Otp             *otpContainer.OtpContainer
	Search          *searchContainer.SearchContainer
	Ads             *adsContainer.AdsContainer
	AdReview        *adReviewContainer.AdReviewContainer
	Setting         *settingContainer.SettingContainer
	Action          *actionContainer.ActionContainer
	Role            *roleContainer.RoleContainer
	RoleAction      *roleActionContainer.RoleActionContainer
	Staff           *staffContainer.StaffContainer
	StaffAuth       *staffAuthContainer.StaffAuthContainer
	StaffRole       *staffRoleContainer.StaffRoleContainer
	StaffGroup      *staffGroupContainer.StaffGroupContainer
}

func NewAppContainer(db *gorm.DB, cfg *config.Config) *AppContainer {
	return &AppContainer{
		Customer:        customerContainer.NewCustomerContainer(db, cfg.DefaultCountryCode),
		CustomerProfile: customerprofileContainer.NewCustomerProfileContainer(db),
		Otp:             otpContainer.NewOtpContainer(cfg.CacheServiceURL, cfg.BrokerURL, cfg.OtpSubject, cfg.DefaultCountryCode, cfg.OtpResendAfter),
		Search:          searchContainer.NewSearchContainer(db, cfg.CacheServiceURL),
		Ads:             adsContainer.NewAdsContainer(db, cfg),
		AdReview:        adReviewContainer.NewAdReviewContainer(db, cfg.CacheServiceURL),
		Setting:         settingContainer.NewSettingContainer(db),
		Action:          actionContainer.NewActionContainer(db),
		Role:            roleContainer.NewRoleContainer(db),
		RoleAction:      roleActionContainer.NewRoleActionContainer(db),
		Staff:           staffContainer.NewStaffContainer(db, cfg.DefaultCountryCode),
		StaffAuth:       staffAuthContainer.NewStaffAuthContainer(db, cfg.DefaultCountryCode),
		StaffRole:       staffRoleContainer.NewStaffRoleContainer(db),
		StaffGroup:      staffGroupContainer.NewStaffGroupContainer(db),
	}
}
