package model

// Party defines ones of the parties involved in a transaction
type Party struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	AccountType       int    `json:"account_type,omitempty"`
	Address           string `json:"address,omitempty"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name,omitempty"`
}

// Valid checks if the given Party are valid or not
func (a *Party) Valid() bool {

	// TODO, depending on the data model, implement this
	return true
}
