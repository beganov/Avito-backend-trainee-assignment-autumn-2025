package errs

import "errors"

type ErrorCode string

const (
	CodeTeamExists  ErrorCode = "TEAM_EXISTS"
	CodePRExists    ErrorCode = "PR_EXISTS"
	CodePRMerged    ErrorCode = "PR_MERGED"
	CodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	CodeNoCandidate ErrorCode = "NO_CANDIDATE"
	CodeNotFound    ErrorCode = "NOT_FOUND"
)

var (
	ErrTeamExists  = errors.New("team already exists")
	ErrPRExists    = errors.New("PR id already exists")
	ErrPRMerged    = errors.New("PR_MERGED")
	ErrNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
	ErrNotFound    = errors.New("resource not found")
)

type ErrorResponse struct {
	Error struct {
		Code    ErrorCode `json:"code"`
		Message string    `json:"message"`
	} `json:"error"`
}

func NewErrorResponse(code ErrorCode, message string) ErrorResponse {
	var resp ErrorResponse
	resp.Error.Code = code
	resp.Error.Message = message
	return resp
}

func TeamExists() ErrorResponse {
	return NewErrorResponse(CodeTeamExists, ErrTeamExists.Error())
}

func PRExists() ErrorResponse {
	return NewErrorResponse(CodePRExists, ErrPRExists.Error())
}

func PRMerged() ErrorResponse {
	return NewErrorResponse(CodePRMerged, ErrPRMerged.Error())
}

func NotAssigned() ErrorResponse {
	return NewErrorResponse(CodeNotAssigned, ErrNotAssigned.Error())
}

func NoCandidate() ErrorResponse {
	return NewErrorResponse(CodeNoCandidate, ErrNoCandidate.Error())
}

func NotFound() ErrorResponse {
	return NewErrorResponse(CodeNotFound, ErrNotFound.Error())
}
