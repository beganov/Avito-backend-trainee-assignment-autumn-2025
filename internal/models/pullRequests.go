package models

// PullRequest represents a full Pull Request entity
// Used for detailed PR
type PullRequest struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"` // OPEN, MERGED
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"createdAt,omitempty"`
	MergedAt          string   `json:"mergedAt,omitempty"`
}

// PullRequestShort represents a simplified view of a Pull Request
// Used for
type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"` // OPEN, MERGED
}

// PRReassign represents the request for reassigning a reviewer
// Used in the reassign reviewer operation
type PRReassign struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

// PRReassignResponse represents the response after successfully reassigning a reviewer
// Used in the reassign reviewer operation
type PRReassignResponse struct {
	PullRequest PullRequest `json:"pr"`
	ReplacedBy  string      `json:"replaced_by"`
}

type PRResponse struct {
	PullRequest PullRequest `json:"pr"`
}
