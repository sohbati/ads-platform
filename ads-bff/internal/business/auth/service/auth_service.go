package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ads-bff/internal/business/auth/model"
	backendclient "ads-bff/internal/core/client/backend"
	cacheclient "ads-bff/internal/core/client/cache"
	"ads-bff/internal/core/config"

	"github.com/google/uuid"
)

type AuthService interface {
	SendOTP(ctx context.Context, mobile string) (int, []byte, error)
	LoginWithOTP(ctx context.Context, mobile, otp string) (*model.LoginResponse, string, int, []byte, error)
	GetCurrentUser(ctx context.Context, sessionID string) (*model.SessionUser, error)
	Logout(ctx context.Context, sessionID string) error
}

type authService struct {
	cfg     *config.Config
	backend *backendclient.Client
	cache   *cacheclient.Client
}

func NewAuthService(cfg *config.Config, backend *backendclient.Client, cache *cacheclient.Client) AuthService {
	return &authService{cfg: cfg, backend: backend, cache: cache}
}

func (s *authService) SendOTP(ctx context.Context, mobile string) (int, []byte, error) {
	return s.backend.SendOTP(ctx, mobile)
}

func (s *authService) LoginWithOTP(ctx context.Context, mobile, otp string) (*model.LoginResponse, string, int, []byte, error) {
	verifyResp, status, body, err := s.backend.VerifyOTP(ctx, mobile, otp)
	if err != nil {
		return nil, "", http.StatusBadGateway, nil, err
	}
	if status != http.StatusOK {
		return nil, "", status, body, nil
	}
	if verifyResp == nil || !verifyResp.Verified {
		return nil, "", http.StatusUnauthorized, body, nil
	}

	user, userStatus, userBody, err := s.backend.GetUserByMobile(ctx, mobile)
	if err != nil {
		return nil, "", http.StatusBadGateway, nil, err
	}
	if userStatus == http.StatusNotFound {
		user, userStatus, userBody, err = s.backend.RegisterUserByMobile(ctx, mobile)
		if err != nil {
			return nil, "", http.StatusBadGateway, nil, err
		}
	}
	if userStatus != http.StatusOK {
		return nil, "", userStatus, userBody, nil
	}
	if user == nil {
		return nil, "", http.StatusBadGateway, nil, fmt.Errorf("user missing after register")
	}

	sessionUser := model.SessionUser{
		ID:         user.ID,
		Name:       user.Name,
		Mobile:     user.Mobile,
		NationalId: user.NationalId,
	}

	sessionData, err := json.Marshal(sessionUser)
	if err != nil {
		return nil, "", http.StatusInternalServerError, nil, err
	}

	sessionID := uuid.NewString()
	cacheKey := sessionCacheKey(sessionID)
	if err := s.cache.StoreSession(ctx, cacheKey, string(sessionData)); err != nil {
		return nil, "", http.StatusBadGateway, nil, fmt.Errorf("store session: %w", err)
	}

	return &model.LoginResponse{
		Authenticated: true,
		User:          sessionUser,
	}, sessionID, http.StatusOK, nil, nil
}

func (s *authService) GetCurrentUser(ctx context.Context, sessionID string) (*model.SessionUser, error) {
	if sessionID == "" {
		return nil, cacheclient.ErrSessionNotFound
	}

	data, err := s.cache.GetSession(ctx, sessionCacheKey(sessionID))
	if err != nil {
		return nil, err
	}

	var user model.SessionUser
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, fmt.Errorf("parse session user: %w", err)
	}
	return &user, nil
}

func (s *authService) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	return s.cache.DeleteSession(ctx, sessionCacheKey(sessionID))
}

func sessionCacheKey(sessionID string) string {
	return "session:" + sessionID
}
