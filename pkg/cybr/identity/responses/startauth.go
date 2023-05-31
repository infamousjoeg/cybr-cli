package responses

import (
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/responses/shared"
)

// StartAuthentication is the response from the start authentication request
type StartAuthentication struct {
	Success         bool          `json:"success"`
	Result          shared.Result `json:"Result"`
	Message         *string       `json:"Message"`
	MessageID       *string       `json:"MessageID"`
	Exception       *string       `json:"Exception"`
	ErrorID         *string       `json:"ErrorID"`
	ErrorCode       *string       `json:"ErrorCode"`
	IsSoftError     bool          `json:"IsSoftError"`
	InnerExceptions *string       `json:"InnerExceptions"`
}
