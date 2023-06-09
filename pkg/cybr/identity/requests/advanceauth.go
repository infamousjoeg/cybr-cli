package requests

// AdvanceAuthentication is the request body for the AdvanceAuthentication API call
type AdvanceAuthentication struct {
	SessionID   string `json:"SessionId"`
	MechanismID string `json:"MechanismId"`
	Action      string `json:"Action"`
	Answer      string `json:"Answer"`
}
