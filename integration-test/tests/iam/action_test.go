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

func TestCreateAndGetAction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	code := uniqueCode("ad.approve")
	status, created, errResp, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/actions", code, "Approve ad")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("create status=%d error=%+v", status, errResp)
	}
	if created.ID <= 0 || created.Code != code || created.Name != "Approve ad" {
		t.Fatalf("created=%+v", created)
	}

	status, got, errResp, err := iamapi.GetNamed(ctx, backURL, "/api/v1/panel/actions/"+strconv.FormatInt(created.ID, 10))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if status != http.StatusOK || got.ID != created.ID || got.Code != code {
		t.Fatalf("get status=%d item=%+v error=%+v", status, got, errResp)
	}

	status, listed, errResp, err := iamapi.ListNamed(ctx, backURL, "/api/v1/panel/actions")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if status != http.StatusOK || !hasCode(listed, code) {
		t.Fatalf("list status=%d missing %s error=%+v", status, code, errResp)
	}
}

func TestUpdateActionName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	code := uniqueCode("ad.reject")
	_, created, _, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/actions", code, "Reject")
	if err != nil {
		t.Fatal(err)
	}

	status, updated, errResp, err := iamapi.PutName(ctx, backURL, "/api/v1/panel/actions/"+strconv.FormatInt(created.ID, 10), "Reject ad")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if status != http.StatusOK || updated.Code != code || updated.Name != "Reject ad" {
		t.Fatalf("update status=%d item=%+v error=%+v", status, updated, errResp)
	}
}

func TestCreateActionRejectsInvalidCode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, _, errResp, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/actions", "Approve Ad", "Approve")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusBadRequest || errResp.Error != "ACTION_INVALID_CODE" {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
}

func TestCreateActionDuplicateCode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	code := uniqueCode("ad.list")
	status, _, errResp, err := iamapi.PostNamed(ctx, backURL, "/api/v1/panel/actions", code, "List")
	if err != nil || status != http.StatusCreated {
		t.Fatalf("first create status=%d error=%+v err=%v", status, errResp, err)
	}

	status, _, errResp, err = iamapi.PostNamed(ctx, backURL, "/api/v1/panel/actions", code, "List again")
	if err != nil {
		t.Fatalf("second create: %v", err)
	}
	if status != http.StatusConflict || errResp.Error != "ACTION_DUPLICATE_CODE" {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
}
