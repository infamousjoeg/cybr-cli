package requests

// GetAccountPassword get an accounts password from PAS. All fields are optional expect if ticketing system is integrated
type GetAccountPassword struct {
	Reason              string `json:"reason,omitempty"`
	TicketingSystemName string `json:"TicketingSystemName,omitempty"`
	TicketID            string `json:"TicketId,omitempty"`
	Version             int    `json:"Version,omitempty"`
	ActionType          string `json:"ActionType,omitempty"`
	IsUse               bool   `json:"isUse,omitempty"`
	Machine             string `json:"Machine,omitempty"`
}
