package errs

type ErrorCode string

const (
	CodeTeamExists  ErrorCode = "TEAM_EXISTS"
	CodePRExists    ErrorCode = "PR_EXISTS"
	CodePRMerged    ErrorCode = "PR_MERGED"
	CodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	CodeNoCandidate ErrorCode = "NO_CANDIDATE"
	CodeNotFound    ErrorCode = "NOT_FOUND"
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
	return NewErrorResponse(CodeTeamExists, "team already exists")
}

func NotFound() ErrorResponse {
	return NewErrorResponse(CodeNotFound, "resource not found")
}

func PRExists() ErrorResponse {
	return NewErrorResponse(CodePRExists, "PR id already exists")
}

func PRMerged() ErrorResponse {
	return NewErrorResponse(CodePRMerged, "cannot reassign on merged PR")
}

func NotAssigned() ErrorResponse {
	return NewErrorResponse(CodeNotAssigned, "reviewer is not assigned to this PR")
}

func NoCandidate() ErrorResponse {
	return NewErrorResponse(CodeNoCandidate, "no active replacement candidate in team")
}
