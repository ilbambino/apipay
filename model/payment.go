package model

// PaymentID is the type of the IDs of payments
type PaymentID string

// Payment defines a payment in the system
// TODO, probably this can be generalized to a Transaction or similar, once more types are added
type Payment struct {
	Type           string     `json:"type"`
	ID             PaymentID  `json:"id"`
	Version        uint       `json:"version"`
	OrganisationID string     `json:"organisation_id"`
	Attributes     Attributes `json:"attributes"`
}

// Valid checks if the given payment is valid or not
func (p *Payment) Valid() bool {

	return p.Type == "Payment" && len(p.ID) > 0 && len(p.OrganisationID) > 0 && p.Attributes.Valid()
}
