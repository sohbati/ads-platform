//go:build integration

package iam

import (
	"context"
	"fmt"
	"testing"
	"time"

	"integration-test/internal/iamapi"

	"github.com/jackc/pgx/v5"
)

func uniqueCode(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func insertUser(t *testing.T, ctx context.Context) int64 {
	t.Helper()
	conn, err := pgx.Connect(ctx, pgHostDSN)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}
	defer conn.Close(ctx)

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	var id int64
	err = conn.QueryRow(ctx, `
		INSERT INTO ads_platform_schema."user" (name, mobile, national_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, "Staff "+suffix, "09"+suffix[len(suffix)-9:], suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	return id
}

func hasCode(items []iamapi.Named, code string) bool {
	for _, item := range items {
		if item.Code == code {
			return true
		}
	}
	return false
}
