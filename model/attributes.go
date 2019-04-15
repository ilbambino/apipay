package model

// SenderCharges contains the information of the charges to the sender
type SenderCharges struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// ChargesInformation holds the information about all charges of the transaction
type ChargesInformation struct {
	BearerCode              string          `json:"bearer_code"`
	SenderCharges           []SenderCharges `json:"sender_charges"`
	ReceiverChargesAmount   string          `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string          `json:"receiver_charges_currency"`
}

// Fx is something I dont' kow about TODO
type Fx struct {
	ContractReference string `json:"contract_reference"`
	ExchangeRate      string `json:"exchange_rate"`
	OriginalAmount    string `json:"original_amount"`
	OriginalCurrency  string `json:"original_currency"`
}

// Attributes holds extra information of the transaction
type Attributes struct {
	Amount               string             `json:"amount"`
	BeneficiaryParty     Party              `json:"beneficiary_party"`
	ChargesInformation   ChargesInformation `json:"charges_information"`
	Currency             string             `json:"currency"`
	DebtorParty          Party              `json:"debtor_party"`
	EndToEndReference    string             `json:"end_to_end_reference"`
	Fx                   Fx                 `json:"fx"`
	NumericReference     string             `json:"numeric_reference"`
	PaymentID            string             `json:"payment_id"`
	PaymentPurpose       string             `json:"payment_purpose"`
	PaymentScheme        string             `json:"payment_scheme"`
	PaymentType          string             `json:"payment_type"`
	ProcessingDate       string             `json:"processing_date"`
	Reference            string             `json:"reference"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type"`
	SchemePaymentType    string             `json:"scheme_payment_type"`
	SponsorParty         Party              `json:"sponsor_party"`
}

// Valid checks if the given attributes are valid or not
func (a *Attributes) Valid() bool {

	// TODO, depending on the data model, implement this
	return true
}
