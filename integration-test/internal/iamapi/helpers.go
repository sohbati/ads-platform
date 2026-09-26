package iamapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Error      string   `json:"error"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params"`
}

type Named struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

const panelRoot = "/api/v1/panel"

func PostNamed(ctx context.Context, backURL, path, code, name string) (int, Named, ErrorResponse, error) {
	req, err := newJSONRequest(ctx, http.MethodPost, backURL+path, map[string]string{"code": code, "name": name})
	if err != nil {
		return 0, Named{}, ErrorResponse{}, err
	}
	var item Named
	status, failure, err := doJSON(req, http.StatusCreated, &item)
	return status, item, failure, err
}

func GetNamed(ctx context.Context, backURL, path string) (int, Named, ErrorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, backURL+path, nil)
	if err != nil {
		return 0, Named{}, ErrorResponse{}, err
	}
	var item Named
	status, failure, err := doJSON(req, http.StatusOK, &item)
	return status, item, failure, err
}

func ListNamed(ctx context.Context, backURL, path string) (int, []Named, ErrorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, backURL+path, nil)
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}
	var items []Named
	status, failure, err := doJSON(req, http.StatusOK, &items)
	if items == nil {
		items = []Named{}
	}
	return status, items, failure, err
}

func PutName(ctx context.Context, backURL, path, name string) (int, Named, ErrorResponse, error) {
	req, err := newJSONRequest(ctx, http.MethodPut, backURL+path, map[string]string{"name": name})
	if err != nil {
		return 0, Named{}, ErrorResponse{}, err
	}
	var item Named
	status, failure, err := doJSON(req, http.StatusOK, &item)
	return status, item, failure, err
}

func Delete(ctx context.Context, backURL, path string) (int, ErrorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, backURL+path, nil)
	if err != nil {
		return 0, ErrorResponse{}, err
	}
	status, failure, err := doJSON(req, http.StatusNoContent, nil)
	return status, failure, err
}

func PutRoleActions(ctx context.Context, backURL string, roleID int64, actionIDs []int64) (int, []Named, ErrorResponse, error) {
	if actionIDs == nil {
		actionIDs = []int64{}
	}
	path := backURL + panelRoot + "/roles/" + strconv.FormatInt(roleID, 10) + "/actions"
	req, err := newJSONRequest(ctx, http.MethodPut, path, map[string]any{"action_ids": actionIDs})
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}
	var items []Named
	status, failure, err := doJSON(req, http.StatusOK, &items)
	if items == nil {
		items = []Named{}
	}
	return status, items, failure, err
}

func PostRoleAction(ctx context.Context, backURL string, roleID, actionID int64) (int, []Named, ErrorResponse, error) {
	path := backURL + panelRoot + "/roles/" + strconv.FormatInt(roleID, 10) + "/actions/" + strconv.FormatInt(actionID, 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}
	var items []Named
	status, failure, err := doJSON(req, http.StatusOK, &items)
	if items == nil {
		items = []Named{}
	}
	return status, items, failure, err
}

func DeleteRoleAction(ctx context.Context, backURL string, roleID, actionID int64) (int, ErrorResponse, error) {
	path := panelRoot + "/roles/" + strconv.FormatInt(roleID, 10) + "/actions/" + strconv.FormatInt(actionID, 10)
	return Delete(ctx, backURL, path)
}

func ListRoleActions(ctx context.Context, backURL string, roleID int64) (int, []Named, ErrorResponse, error) {
	return ListNamed(ctx, backURL, panelRoot+"/roles/"+strconv.FormatInt(roleID, 10)+"/actions")
}

func PostUserRole(ctx context.Context, backURL string, userID, roleID int64) (int, []Named, ErrorResponse, error) {
	path := backURL + panelRoot + "/users/" + strconv.FormatInt(userID, 10) + "/roles/" + strconv.FormatInt(roleID, 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}
	var items []Named
	status, failure, err := doJSON(req, http.StatusOK, &items)
	if items == nil {
		items = []Named{}
	}
	return status, items, failure, err
}

func DeleteUserRole(ctx context.Context, backURL string, userID, roleID int64) (int, ErrorResponse, error) {
	path := panelRoot + "/users/" + strconv.FormatInt(userID, 10) + "/roles/" + strconv.FormatInt(roleID, 10)
	return Delete(ctx, backURL, path)
}

func ListUserRoles(ctx context.Context, backURL string, userID int64) (int, []Named, ErrorResponse, error) {
	return ListNamed(ctx, backURL, panelRoot+"/users/"+strconv.FormatInt(userID, 10)+"/roles")
}

func ListUserActions(ctx context.Context, backURL string, userID int64) (int, []Named, ErrorResponse, error) {
	return ListNamed(ctx, backURL, panelRoot+"/users/"+strconv.FormatInt(userID, 10)+"/actions")
}

func ListUserGroups(ctx context.Context, backURL string, userID int64) (int, []Named, ErrorResponse, error) {
	return ListNamed(ctx, backURL, panelRoot+"/users/"+strconv.FormatInt(userID, 10)+"/groups")
}

func PostGroupUser(ctx context.Context, backURL string, groupID, userID int64) (int, ErrorResponse, error) {
	path := backURL + panelRoot + "/user-groups/" + strconv.FormatInt(groupID, 10) + "/users/" + strconv.FormatInt(userID, 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
	if err != nil {
		return 0, ErrorResponse{}, err
	}
	status, failure, err := doJSON(req, http.StatusNoContent, nil)
	return status, failure, err
}

func ListGroupUsers(ctx context.Context, backURL string, groupID int64) (int, []int64, ErrorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, backURL+panelRoot+"/user-groups/"+strconv.FormatInt(groupID, 10)+"/users", nil)
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}
	var body struct {
		UserIDs []int64 `json:"user_ids"`
	}
	status, failure, err := doJSON(req, http.StatusOK, &body)
	if body.UserIDs == nil {
		body.UserIDs = []int64{}
	}
	return status, body.UserIDs, failure, err
}

func newJSONRequest(ctx context.Context, method, url string, payload any) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func doJSON(req *http.Request, okStatus int, dest any) (int, ErrorResponse, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, ErrorResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, ErrorResponse{}, err
	}

	if resp.StatusCode != okStatus {
		var failure ErrorResponse
		_ = json.Unmarshal(body, &failure)
		return resp.StatusCode, failure, nil
	}

	if dest != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dest); err != nil {
			return resp.StatusCode, ErrorResponse{}, fmt.Errorf("decode: %w; body=%s", err, string(body))
		}
	}
	return resp.StatusCode, ErrorResponse{}, nil
}
