//go:build integration

package iam

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/iamapi"
)

func TestUserRoleGrantsActions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := insertUser(t, ctx)
	_, action, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/actions", uniqueCode("ad.approve"), "Approve")
	if err != nil {
		t.Fatal(err)
	}
	_, role, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/roles", uniqueCode("moderator"), "Moderator")
	if err != nil {
		t.Fatal(err)
	}
	status, _, errResp, err := iamapi.PutRoleActions(ctx, backURL, role.ID, []int64{action.ID})
	if err != nil || status != http.StatusOK {
		t.Fatalf("role actions status=%d error=%+v err=%v", status, errResp, err)
	}

	status, actions, errResp, err := iamapi.ListUserActions(ctx, backURL, userID)
	if err != nil {
		t.Fatalf("actions before role: %v", err)
	}
	if status != http.StatusOK || hasCode(actions, action.Code) {
		t.Fatalf("expected no actions before user_role, status=%d items=%+v error=%+v", status, actions, errResp)
	}

	status, roles, errResp, err := iamapi.PostUserRole(ctx, backURL, userID, role.ID)
	if err != nil || status != http.StatusOK {
		t.Fatalf("assign role status=%d error=%+v err=%v", status, errResp, err)
	}
	if !hasCode(roles, role.Code) {
		t.Fatalf("user roles=%+v", roles)
	}

	status, actions, errResp, err = iamapi.ListUserActions(ctx, backURL, userID)
	if err != nil {
		t.Fatalf("actions after role: %v", err)
	}
	if status != http.StatusOK || !hasCode(actions, action.Code) {
		t.Fatalf("expected granted action, status=%d items=%+v error=%+v", status, actions, errResp)
	}

	status, errResp, err = iamapi.DeleteUserRole(ctx, backURL, userID, role.ID)
	if err != nil || status != http.StatusNoContent {
		t.Fatalf("remove role status=%d error=%+v err=%v", status, errResp, err)
	}
	_, actions, _, err = iamapi.ListUserActions(ctx, backURL, userID)
	if err != nil {
		t.Fatal(err)
	}
	if hasCode(actions, action.Code) {
		t.Fatalf("action remained after removing user_role: %+v", actions)
	}
}

func TestUnknownUserActions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, _, errResp, err := iamapi.ListUserActions(ctx, backURL, 9_999_999)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusNotFound || errResp.Error != "USER_NOT_FOUND" {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
}
