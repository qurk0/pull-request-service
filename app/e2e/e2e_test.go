package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/qurk0/pr-service/internal/api/dto"
)

const baseUrl = "http://localhost:8080"

var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}

type addTeamResponse struct {
	Team dto.TeamResponse `json:"team"`
}

type prResponse struct {
	Pr dto.PR `json:"pr"`
}

func TestPRLifecycleE2E(t *testing.T) {
	suffix := time.Now().UnixNano()
	teamName := fmt.Sprintf("backend-%d", suffix)

	userAuthor := fmt.Sprintf("u-author-%d", suffix)
	userRev1 := fmt.Sprintf("u-rev1-%d", suffix)
	userRev2 := fmt.Sprintf("u-rev2-%d", suffix)
	userRev3 := fmt.Sprintf("u-rev3-%d", suffix)
	userRev4 := fmt.Sprintf("u-rev4-%d", suffix)
	userRev5 := fmt.Sprintf("u-rev5-%d", suffix)
	userRev6 := fmt.Sprintf("u-rev6-%d", suffix)

	prID := fmt.Sprintf("pr-%d", suffix)

	t.Log("step 1: /team/add")
	addReq := dto.AddTeamRequest{
		TeamName: teamName,
		Members: []dto.Member{
			{Id: userAuthor, Username: "author", IsActive: true},
			{Id: userRev1, Username: "rev1", IsActive: true},
			{Id: userRev2, Username: "rev2", IsActive: true},
			{Id: userRev3, Username: "rev3", IsActive: true},
			{Id: userRev4, Username: "rev4", IsActive: true},
			{Id: userRev5, Username: "rev5", IsActive: true},
			{Id: userRev6, Username: "rev6", IsActive: true},
		},
	}

	var addResp addTeamResponse
	postJSON(t, "/team/add", addReq, http.StatusCreated, &addResp)

	if addResp.Team.TeamName != teamName {
		t.Fatalf("expected team_name %q, got %q", teamName, addReq.TeamName)
	}

	if len(addResp.Team.Members) != 7 {
		t.Fatalf("expected 7 members, got %d", len(addResp.Team.Members))
	}

	t.Log("step 2: /pullRequest/create")
	createReq := dto.CreatePRRequest{
		PRID:     prID,
		PRNamme:  "Test PR",
		AuthorID: userAuthor,
	}

	var createResp prResponse
	postJSON(t, "/pullRequest/create", createReq, http.StatusCreated, &createResp)

	pr := createResp.Pr
	if pr.PRID != prID {
		t.Fatalf("expected pr_id = %q, got %q", prID, pr.PRID)
	}
	if pr.AuthorID != userAuthor {
		t.Fatalf("expected author_id %q, got %q", userAuthor, pr.AuthorID)
	}
	if pr.Status != "OPEN" {
		t.Fatalf("expected status OPEN, got %s", pr.Status)
	}
	if len(pr.AssignedReviewers) == 0 {
		t.Fatal("expected at least 1 assigned reviewer, got 0")
	}
	for _, r := range pr.AssignedReviewers {
		if r == userAuthor {
			t.Fatalf("author %q must not be in assigned reviewers", userAuthor)
		}
	}

	t.Log("step3: /users/getReview")
	reviewer := pr.AssignedReviewers[0]

	var reviewResp dto.GetReviewResponse
	getJSON(t, fmt.Sprintf("/users/getReview?user_id=%s", reviewer), http.StatusOK, &reviewResp)

	if reviewResp.UserID != reviewer {
		t.Fatalf("expected user_id %q in getReview, got %q", reviewer, reviewResp.UserID)
	}

	found := false
	for _, short := range reviewResp.RequestsShort {
		if short.ID == prID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected PR %q to be in getReview list for %q", prID, reviewer)
	}

	t.Log("step 4: /pullRequest/reassign")
	reassignReq := dto.ReassignPRRequest{
		PRID:          prID,
		OldReviewerID: reviewer,
	}

	var reassignResp prResponse
	postJSON(t, "/pullRequest/reassign", reassignReq, http.StatusOK, &reassignResp)

	prAfterReassign := reassignResp.Pr
	if prAfterReassign.Status != "OPEN" {
		t.Fatalf("expected status OPEN after reassign, got %s", prAfterReassign.Status)
	}

	if len(prAfterReassign.AssignedReviewers) == 0 {
		t.Fatalf("expected at least 1 reviewer after reassign, got 0")
	}
	for _, r := range prAfterReassign.AssignedReviewers {
		if r == reviewer {
			t.Fatalf("old reviewer %q must not be in assigned_reviewers after reassign", reviewer)
		}

		if r == userAuthor {
			t.Fatalf("author %q mut not be a reviewer", userAuthor)
		}
	}

	t.Log("step 5: /pullRequest/merge")
	mergeReq := dto.PRMergeRequest{
		PRID: prID,
	}

	var mergeResp prResponse
	postJSON(t, "/pullRequest/merge", mergeReq, http.StatusOK, &mergeResp)

	prAfterMerge := mergeResp.Pr
	if prAfterMerge.Status != "MERGED" {
		t.Fatalf("expected status MERGED after merge, got %q", prAfterMerge.Status)
	}
	if prAfterMerge.MergedAt == nil {
		t.Fatal("expected mergedAt to be non-nil value after merge")
	}

	t.Log("step 6: /pullRequest/reassign on merged PR (expect 409 PR_MERGED)")

	old := prAfterReassign.AssignedReviewers[0]

	rawResp := postJSON(t, "/pullRequest/reassign", dto.ReassignPRRequest{
		PRID:          prID,
		OldReviewerID: old,
	}, http.StatusConflict, nil)

	var errResp dto.ErrorResponse
	if err := json.NewDecoder(rawResp.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Error.Code != dto.ErrCodePrMerged {
		t.Fatalf("expected error code %q, got %q", dto.ErrCodePrMerged, errResp.Error.Code)
	}
}

func postJSON(t *testing.T, path string, body any, expectedStatus int, respTarget any) *http.Response {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("failed to encode request body: %v", err)
		}
	}

	req, err := http.NewRequest(http.MethodPost, baseUrl+path, &buf)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {

		t.Fatalf("request %s %s failed: %v", http.MethodPost, path, err)
	}

	if resp.StatusCode != expectedStatus {
		_ = resp.Body.Close()
	}

	if respTarget != nil {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(respTarget); err != nil {
			t.Fatalf("failed to decode response body for %s %s: %v", http.MethodPost, path, err)
		}

	}

	return resp
}

func getJSON(t *testing.T, path string, expectedStatus int, respTarget any) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, baseUrl+path, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request %s failed: %v", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Fatalf("expected status %d for GET %s, got %d", expectedStatus, path, resp.StatusCode)
	}

	if respTarget != nil {
		if err := json.NewDecoder(resp.Body).Decode(respTarget); err != nil {
			t.Fatalf("failed to decode response body for GET %s: %v", path, err)
		}
	}
}
