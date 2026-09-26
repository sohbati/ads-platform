//go:build integration

package iam

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"integration-test/internal/iamapi"
)

func TestCreateAndGetRole(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	code := uniqueCode("moderator")
	status, created, errResp, err := iamapi.PostNamed(ctx, backURL, "/api/v1/roles", code, "Moderator")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if status != http.StatusCreated || created.Code != code {
		t.Fatalf("create status=%d item=%+v error=%+v", status, created, errResp)
	}

	status, got, errResp, err := iamapi.GetNamed(ctx, backURL, "/api/v1/roles/"+strconv.FormatInt(created.ID, 10))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if status != http.StatusOK || got.ID != created.ID {
		t.Fatalf("get status=%d item=%+v error=%+v", status, got, errResp)
	}
}

func TestCreateRoleRejectsInvalidCode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, _, errResp, err := iamapi.PostNamed(ctx, backURL, "/api/v1/roles", "Mod Role", "Moderator")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusBadRequest || errResp.Error != "ROLE_INVALID_CODE" {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
}
