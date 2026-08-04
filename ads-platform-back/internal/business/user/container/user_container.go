package container

// Container acts as the application's dependency container (composition root).
// Its responsibility is to initialize and wire together all core components
// such as repositories, services, and handlers.
//
// The container centralizes dependency creation so that other parts of the
// application (like the router) do not need to know how objects are constructed.
// Instead, they simply receive ready-to-use instances.
//
// This approach provides several benefits:
// - Keeps object creation in one place
// - Promotes clean separation of concerns
// - Makes testing easier by allowing dependencies to be replaced with mocks
// - Prevents tight coupling between layers
//
// In general, the dependency flow in the application is:
//
//   Handler → Service → Repository
//
// The container is responsible for constructing these layers and injecting
// them into each other during application startup.

import (
	"ads-platform/internal/business/user/handler"
	repoimpl "ads-platform/internal/business/user/repository/impl"
	serviceimpl "ads-platform/internal/business/user/service/impl"

	"gorm.io/gorm"
)

type UserContainer struct {
	UserHandler *handler.UserHandler
}

func NewUserContainer(db *gorm.DB) *UserContainer {

	// repositories
	userRepo := repoimpl.NewUserRepository(db)

	// services
	userService := serviceimpl.NewUserService(userRepo)

	// handlers
	userHandler := handler.NewUserHandler(userService)

	return &UserContainer{
		UserHandler: userHandler,
	}
}
