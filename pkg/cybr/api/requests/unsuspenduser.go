package requests

// UnsuspendUser request used when unsuspending user
type UnsuspendUser struct {
	Suspended bool `json:"Suspended"`
}
