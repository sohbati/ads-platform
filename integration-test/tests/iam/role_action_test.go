//go:build integration

package iam

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/iamapi"
)

func TestAssignActionsToRole(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, approve, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/actions", uniqueCode("ad.approve"), "Approve")
	if err != nil {
		t.Fatal(err)
	}
	_, reject, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/actions", uniqueCode("ad.reject"), "Reject")
	if err != nil {
		t.Fatal(err)
	}
	_, role, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/roles", uniqueCode("moderator"), "Moderator")
	if err != nil {
		t.Fatal(err)
	}

	status, items, errResp, err := iamapi.PutRoleActions(ctx, backURL, role.ID, []int64{approve.ID, reject.ID})
	if err != nil {
		t.Fatalf("set: %v", err)
	}
	if status != http.StatusOK || len(items) != 2 {
		t.Fatalf("set status=%d items=%+v error=%+v", status, items, errResp)
	}

	status, listed, errResp, err := iamapi.ListRoleActions(ctx, backURL, role.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if status != http.StatusOK || !hasCode(listed, approve.Code) || !hasCode(listed, reject.Code) {
		t.Fatalf("list status=%d items=%+v error=%+v", status, listed, errResp)
	}

	status, errResp, err = iamapi.DeleteRoleAction(ctx, backURL, role.ID, reject.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if status != http.StatusNoContent {
		t.Fatalf("delete status=%d error=%+v", status, errResp)
	}

	_, listed, _, err = iamapi.ListRoleActions(ctx, backURL, role.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCode(listed, approve.Code) || hasCode(listed, reject.Code) {
		t.Fatalf("after delete items=%+v", listed)
	}
}

func TestAddRoleActionIdempotent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, action, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/actions", uniqueCode("staff.read"), "Read staff")
	if err != nil {
		t.Fatal(err)
	}
	_, role, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/roles", uniqueCode("support"), "Support")
	if err != nil {
		t.Fatal(err)
	}

	status, items, errResp, err := iamapi.PostRoleAction(ctx, backURL, role.ID, action.ID)
	if err != nil || status != http.StatusOK {
		t.Fatalf("first add status=%d error=%+v err=%v", status, errResp, err)
	}
	status, items, errResp, err = iamapi.PostRoleAction(ctx, backURL, role.ID, action.ID)
	if err != nil || status != http.StatusOK {
		t.Fatalf("second add status=%d error=%+v err=%v", status, errResp, err)
	}
	if len(items) != 1 || items[0].ID != action.ID {
		t.Fatalf("items=%+v", items)
	}
}
