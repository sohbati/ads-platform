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
	GetCurrentCustomer(ctx context.Context, sessionID string) (*model.SessionCustomer, error)
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

	customer, customerStatus, customerBody, err := s.backend.GetCustomerByMobile(ctx, mobile)
	if err != nil {
		return nil, "", http.StatusBadGateway, nil, err
	}
	if customerStatus == http.StatusNotFound {
		customer, customerStatus, customerBody, err = s.backend.RegisterCustomerByMobile(ctx, mobile)
		if err != nil {
			return nil, "", http.StatusBadGateway, nil, err
		}
	}
	if customerStatus != http.StatusOK {
		return nil, "", customerStatus, customerBody, nil
	}
	if customer == nil {
		return nil, "", http.StatusBadGateway, nil, fmt.Errorf("customer missing after register")
	}

	sessionCustomer := model.SessionCustomer{
		ID:         customer.ID,
		Name:       customer.Name,
		Mobile:     customer.Mobile,
		NationalId: customer.NationalId,
	}

	sessionData, err := json.Marshal(sessionCustomer)
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
		Customer:      sessionCustomer,
	}, sessionID, http.StatusOK, nil, nil
}

func (s *authService) GetCurrentCustomer(ctx context.Context, sessionID string) (*model.SessionCustomer, error) {
	if sessionID == "" {
		return nil, cacheclient.ErrSessionNotFound
	}

	data, err := s.cache.GetSession(ctx, sessionCacheKey(sessionID))
	if err != nil {
		return nil, err
	}

	var customer model.SessionCustomer
	if err := json.Unmarshal([]byte(data), &customer); err != nil {
		return nil, fmt.Errorf("parse session customer: %w", err)
	}
	return &customer, nil
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
