package shared

// Challenge contains the challenge response including the mechanisms
type Challenge struct {
	Mechanisms []Mechanism `json:"Mechanisms"`
}

// Mechanism contains the mechanism response
type Mechanism struct {
	AnswerType           string `json:"AnswerType"`
	Name                 string `json:"Name"`
	PromptMechChosen     string `json:"PromptMechChosen"`
	PromptSelectMech     string `json:"PromptSelectMech"`
	MechanismID          string `json:"MechanismId"`
	Enrolled             bool   `json:"Enrolled"`
	MaskedEmailAddress   string `json:"MaskedEmailAddress,omitempty"`
	PartialAddress       string `json:"PartialAddress,omitempty"`
	PartialDeviceAddress string `json:"PartialDeviceAddress,omitempty"`
}
