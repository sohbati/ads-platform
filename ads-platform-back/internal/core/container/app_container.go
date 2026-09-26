package container

import (
	adsContainer "ads-platform/internal/business/ads/container"
	actionContainer "ads-platform/internal/business/iam/action/container"
	roleContainer "ads-platform/internal/business/iam/role/container"
	roleActionContainer "ads-platform/internal/business/iam/role_action/container"
	userContainer "ads-platform/internal/business/iam/user/container"
	userGroupContainer "ads-platform/internal/business/iam/user_group/container"
	userRoleContainer "ads-platform/internal/business/iam/user_role/container"
	otpContainer "ads-platform/internal/business/otp/container"
	searchContainer "ads-platform/internal/business/search/container"
	userprofileContainer "ads-platform/internal/business/userprofile/container"
	"ads-platform/internal/core/config"

	"gorm.io/gorm"
)

type AppContainer struct {
	User        *userContainer.UserContainer
	UserProfile *userprofileContainer.UserProfileContainer
	Otp         *otpContainer.OtpContainer
	Search      *searchContainer.SearchContainer
	Ads         *adsContainer.AdsContainer
	Action      *actionContainer.ActionContainer
	Role        *roleContainer.RoleContainer
	RoleAction  *roleActionContainer.RoleActionContainer
	UserRole    *userRoleContainer.UserRoleContainer
	UserGroup   *userGroupContainer.UserGroupContainer
}

func NewAppContainer(db *gorm.DB, cfg *config.Config) *AppContainer {
	return &AppContainer{
		User:        userContainer.NewUserContainer(db, cfg.DefaultCountryCode),
		UserProfile: userprofileContainer.NewUserProfileContainer(db),
		Otp:         otpContainer.NewOtpContainer(cfg.CacheServiceURL, cfg.BrokerURL, cfg.OtpSubject, cfg.DefaultCountryCode, cfg.OtpResendAfter),
		Search:      searchContainer.NewSearchContainer(db, cfg.CacheServiceURL),
		Ads:         adsContainer.NewAdsContainer(db, cfg),
		Action:      actionContainer.NewActionContainer(db),
		Role:        roleContainer.NewRoleContainer(db),
		RoleAction:  roleActionContainer.NewRoleActionContainer(db),
		UserRole:    userRoleContainer.NewUserRoleContainer(db),
		UserGroup:   userGroupContainer.NewUserGroupContainer(db),
	}
}
