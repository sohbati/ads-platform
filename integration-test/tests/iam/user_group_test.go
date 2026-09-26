//go:build integration

package iam

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/iamapi"
)

func TestUserGroupCategorizesUsers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := insertUser(t, ctx)
	status, group, errResp, err := iamapi.PostNamed(ctx, backURL, "/api/v1/user-groups", uniqueCode("tehran-moderators"), "Tehran")
	if err != nil || status != http.StatusCreated {
		t.Fatalf("create group status=%d error=%+v err=%v", status, errResp, err)
	}

	status, errResp, err = iamapi.PostGroupUser(ctx, backURL, group.ID, userID)
	if err != nil || status != http.StatusNoContent {
		t.Fatalf("add member status=%d error=%+v err=%v", status, errResp, err)
	}

	status, ids, errResp, err := iamapi.ListGroupUsers(ctx, backURL, group.ID)
	if err != nil {
		t.Fatalf("list members: %v", err)
	}
	if status != http.StatusOK || len(ids) != 1 || ids[0] != userID {
		t.Fatalf("members status=%d ids=%v error=%+v", status, ids, errResp)
	}

	status, groups, errResp, err := iamapi.ListUserGroups(ctx, backURL, userID)
	if err != nil {
		t.Fatalf("user groups: %v", err)
	}
	if status != http.StatusOK || !hasCode(groups, group.Code) {
		t.Fatalf("user groups status=%d items=%+v error=%+v", status, groups, errResp)
	}
}

func TestUserGroupDoesNotGrantActions(t *testing.T) {
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
	if status, _, errResp, err := iamapi.PutRoleActions(ctx, backURL, role.ID, []int64{action.ID}); err != nil || status != http.StatusOK {
		t.Fatalf("role actions status=%d error=%+v err=%v", status, errResp, err)
	}

	_, group, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/user-groups", uniqueCode("tehran-moderators"), "Tehran")
	if err != nil {
		t.Fatal(err)
	}
	if status, errResp, err := iamapi.PostGroupUser(ctx, backURL, group.ID, userID); err != nil || status != http.StatusNoContent {
		t.Fatalf("add member status=%d error=%+v err=%v", status, errResp, err)
	}

	status, actions, errResp, err := iamapi.ListUserActions(ctx, backURL, userID)
	if err != nil {
		t.Fatalf("user actions: %v", err)
	}
	if status != http.StatusOK || hasCode(actions, action.Code) {
		t.Fatalf("group must not grant actions, status=%d items=%+v error=%+v", status, actions, errResp)
	}
}
